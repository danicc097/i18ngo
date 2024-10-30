package i18ngo

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/tools/imports"
	"mvdan.cc/gofumpt/format"

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
	// gotempl
	/* var buf bytes.Buffer
	 component := templates.TranslationCode(data)
	err := component.Render(context.Background(), &buf)
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %w", err)
	} */

	src, err := generateWithGoTemplate(data)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func generateWithGoTemplate(data *templates.TemplateData) ([]byte, error) {
	funcMap := template.FuncMap{
		"camelCase": func(s string) string {
			return snaker.ForceLowerCamelIdentifier(s)
		},
		"pascalCase": func(s string) string {
			return snaker.ForceCamelIdentifier(s)
		},
	}

	var tplFsys fs.FS = templateFS
	tmpl := template.Must(template.New("template.go.tpl").Funcs(funcMap).ParseFS(tplFsys, "templates/template.go.tpl"))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte{}, fmt.Errorf("error executing template: %w", err)
	}

	source := buf.Bytes()

	source, err := format.Source(source, format.Options{})
	if err != nil {
		return []byte{}, fmt.Errorf("error formatting generated Go code: %w", err)
	}

	source, err = imports.Process("", source, &imports.Options{
		FormatOnly: true,
		Comments:   true,
	})
	if err != nil {
		return []byte{}, fmt.Errorf("error applying gofumpt to generated Go code: %w", err)
	}

	return source, nil
}

var tmplVarRe = regexp.MustCompile(`{{\s*\.\s*([\w\.]+)\s*}}`)

func extractTemplateVariables(tplStr string) ([]string, error) {
	matches := tmplVarRe.FindAllStringSubmatch(tplStr, -1)
	var vars []string
	for _, match := range matches {
		vars = append(vars, match[1])
	}
	vars = unique(vars)
	return vars, nil
}

func unique[T comparable](input []T) []T {
	m := make(map[T]bool)
	var result []T
	for _, s := range input {
		if !m[s] {
			m[s] = true
			result = append(result, s)
		}
	}
	return result
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

			varnames, err := extractTemplateVariables(msg.Template)
			if err != nil {
				return nil, fmt.Errorf("error extracting template variables: %w", err)
			}
			varsm := map[string]templates.VarData{}
			for _, v := range varnames {
				varsm[v] = templates.VarData{
					Name:  v,
					Type:  "interface{}",
					Param: snaker.ForceLowerCamelIdentifier(v),
				}
			}

			// TODO: allow custom imports --> enables e.g. User.Username, User.Gender, etc.
			// in the future for easier custom_templates.
			exprVars := make([]string, 0, len(msg.Variables))
			for name, typ := range msg.Variables {
				varsm[name] = templates.VarData{
					Name:  name,
					Type:  typ,
					Param: snaker.ForceLowerCamelIdentifier(name),
				}
				exprVars = append(exprVars, snaker.ForceLowerCamelIdentifier(name))
			}
			vars := make([]templates.VarData, 0, len(varsm))
			for _, v := range varsm {
				vars = append(vars, v)
			}
			sort.Slice(vars, func(i, j int) bool {
				return vars[i].Name < vars[j].Name
			})

			if err := validator.ValidateTemplate(msg.Template, varnames); err != nil {
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
