package i18ngo

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"golang.org/x/text/language"

	"github.com/danicc097/i18ngo/templates"
	"github.com/danicc097/i18ngo/validator"
	"github.com/kenshaw/snaker"
	"gopkg.in/yaml.v3"
)

type LanguageLoader struct {
	translations map[string]templates.Translations
}

func NewLanguageLoader(fsys fs.FS, path string) (*LanguageLoader, error) {
	loader := &LanguageLoader{translations: make(map[string]templates.Translations)}

	if err := validator.ValidateTranslationFiles(fsys, path); err != nil {
		return nil, fmt.Errorf("error validating translation files: %w", err)
	}
	err := fs.WalkDir(fsys, path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(p, ".i18ngo.yaml") {
			file, err := fs.ReadFile(fsys, p)
			if err != nil {
				return err
			}
			var t templates.Translations
			if err := yaml.Unmarshal(file, &t); err != nil {
				return err
			}
			tlFile := p[strings.LastIndex(p, "/")+1:]
			lang := strings.Split(tlFile, ".i18ngo.yaml")[0]
			if _, err = language.Parse(lang); err != nil {
				return fmt.Errorf("invalid locale %s: %w", lang, err)
			}
			loader.translations[lang] = t
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return loader, nil
}

//go:embed templates/template.go.tpl
var templateFS embed.FS

type GenerateOption func(*generateOptions)

type generateOptions struct {
	WithCustomTemplate bool
}

func WithFilesystemTemplate() GenerateOption {
	return func(opts *generateOptions) {
		opts.WithCustomTemplate = true
	}
}

func Generate(data *templates.TemplateData) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data must be non-nil")
	}
	var buf bytes.Buffer
	component := templates.TranslationCode(data)
	err := component.Render(context.Background(), &buf)
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %w", err)
	}

	return buf.Bytes(), nil
}

// GetTranslationData retrieves data for translations in the given path in the filesystem.
// Assumes the filesystem contains a templates/template.go.tpl file to generate from.
// You may extend the default template as desired.
func GetTranslationData(fsys fs.FS, path, pkgName string, opts ...GenerateOption) (*templates.TemplateData, error) {
	optsMap := &generateOptions{}
	for _, o := range opts {
		o(optsMap)
	}

	loader, err := NewLanguageLoader(fsys, path)
	if err != nil {
		return nil, err
	}

	data := templates.TemplateData{
		PkgName:      pkgName,
		Messages:     make([]templates.MessageData, 0),
		Translations: make([]templates.TranslationData, 0),
		Langs:        make([]templates.LangData, 0),
	}

	langKeys := make([]string, 0, len(loader.translations))
	for lang := range loader.translations {
		langKeys = append(langKeys, lang)
	}
	sort.Strings(langKeys)

	for _, lang := range langKeys {
		translations := loader.translations[lang]
		camelLang := snaker.SnakeToCamel(lang)
		data.Langs = append(data.Langs, templates.LangData{CamelLang: camelLang, Lang: lang})

		transData := templates.TranslationData{CamelLang: camelLang}

		msgIDs := make([]string, 0, len(translations.Messages))
		for msgID := range translations.Messages {
			msgIDs = append(msgIDs, msgID)
		}
		sort.Strings(msgIDs)

		for _, msgID := range msgIDs {
			msg := translations.Messages[msgID]
			methodName := snaker.SnakeToCamel(msgID)
			vars := []templates.VarData{}

			// TODO: generate variables with interface{} type instead
			// if var not defined.
			// Also allow custom imports --> enables e.g. User.Username, User.Gender, etc.
			// in the future for easier custom_templates.
			tplVars := make([]string, 0, len(msg.Variables))
			exprVars := make([]string, 0, len(msg.Variables))
			for name, typ := range msg.Variables {
				vars = append(vars, templates.VarData{
					Name:  name,
					Type:  typ,
					Param: snaker.ForceLowerCamelIdentifier(name),
				})
				tplVars = append(tplVars, name)
				exprVars = append(exprVars, snaker.ForceLowerCamelIdentifier(name))
			}
			sort.Slice(vars, func(i, j int) bool {
				return vars[i].Name < vars[j].Name
			})

			if err := validator.ValidateTemplate(msg.Template, tplVars); err != nil {
				return nil, fmt.Errorf("error validating template %q: %w", msg.Template, err)
			}

			for _, tpl := range msg.CustomTemplates {
				if err := validator.ValidateCustomExpression(tpl.Expression, exprVars); err != nil {
					return nil, fmt.Errorf("error validating custom template expression %q: %w", tpl.Expression, err)
				}
			}

			args := ""
			for _, v := range vars {
				args += fmt.Sprintf("%s %s, ", v.Param, v.Type)
			}
			args = strings.TrimSuffix(args, ", ")

			transData.Messages = append(transData.Messages, templates.MessageData{
				CamelLang:       camelLang,
				MethodName:      methodName,
				Args:            args,
				Vars:            vars,
				Template:        msg.Template,
				CustomTemplates: msg.CustomTemplates,
			})
		}
		data.Translations = append(data.Translations, transData)
	}

	data.Messages = data.Translations[0].Messages // all translations have the same messages

	return &data, nil
}
