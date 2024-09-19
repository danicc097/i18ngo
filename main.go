package i18ngo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"sort"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/tools/imports"
	"mvdan.cc/gofumpt/format"

	"github.com/danicc097/i18ngo/validator"
	"github.com/kenshaw/snaker"
	"gopkg.in/yaml.v3"
)

// Message defines the structure for each message entry.
type Message struct {
	Template        string            `yaml:"template"`
	Variables       map[string]string `yaml:"variables"`
	CustomTemplates map[string]string `yaml:"custom_templates"`
}

// Translations holds messages.
type Translations struct {
	Messages map[string]Message `yaml:"messages"`
}

// LanguageLoader will load language files based on the current context language.
type LanguageLoader struct {
	translations map[string]Translations
}

type LangCtxKey struct{}

// I18n struct represents the translation service.
type I18n struct {
	Ctx    context.Context
	Loader *LanguageLoader
}

// NewLanguageLoader initializes and loads YAML-based translations in the given directory.
func NewLanguageLoader(fsys fs.FS, path string) (*LanguageLoader, error) {
	loader := &LanguageLoader{translations: make(map[string]Translations)}

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
			var t Translations
			if err := yaml.Unmarshal(file, &t); err != nil {
				return err
			}
			tlFile := p[strings.LastIndex(p, "/")+1:]
			lang := strings.Split(tlFile, ".i18ngo.yaml")[0]
			fmt.Printf("lang: %v\n", lang)
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

// Generate will generate go code based on the translation files in the given directory.
func Generate(fsys fs.FS, path string) ([]byte, error) {
	loader, err := NewLanguageLoader(fsys, path)
	if err != nil {
		return []byte{}, err
	}
	fmt.Printf("loader.translations: %v\n", loader.translations)

	var buf bytes.Buffer
	buf.WriteString("package translation\n\n")
	buf.WriteString("import (\n\t\"context\"\n\t\"github.com/danicc097/i18ngo\"\n)\n\n")
	buf.WriteString("type T struct {\n\tCtx context.Context\n\tl   *i18ngo.LanguageLoader\n}\n\n")
	buf.WriteString("// New returns a new i18n translator.\n")
	buf.WriteString("func New(l *i18ngo.LanguageLoader) *T {\n\treturn &T{\n\t\tl: l,\n\t}\n}\n\n")

	langs := make([]string, 0, len(loader.translations))
	for k := range loader.translations {
		langs = append(langs, k)
	}
	lang := langs[0]
	messageIDs := make([]string, 0, len(loader.translations[lang].Messages))
	for id := range loader.translations[lang].Messages {
		messageIDs = append(messageIDs, id)
	}
	sort.Strings(messageIDs)

	for _, msgID := range messageIDs {
		msg := loader.translations[lang].Messages[msgID]
		methName := snaker.SnakeToCamel(msgID)

		def := fmt.Sprintf("func (t *T) %s(", methName)

		varNames := make([]string, 0, len(msg.Variables))
		for name := range msg.Variables {
			varNames = append(varNames, name)
		}
		sort.Strings(varNames)

		for _, name := range varNames {
			typ := msg.Variables[name]
			def += fmt.Sprintf("%s %s, ", snaker.ForceLowerCamelIdentifier(name), typ)
		}
		def = strings.TrimSuffix(def, ", ") + ") (string, error) {\n"

		def += "\tdata := struct{\n"
		for _, name := range varNames {
			typ := msg.Variables[name]
			def += fmt.Sprintf("\t\t%s %s\n", snaker.SnakeToCamel(name), typ)
		}
		def += "\t\tI18ngoCustomTemplateExpr string\n"
		def += "\t} {\n"
		for _, name := range varNames {
			def += fmt.Sprintf("\t\t%s: %s,\n", snaker.SnakeToCamel(name), snaker.ForceLowerCamelIdentifier(name))
		}
		def += "\t}\n"

		if len(msg.CustomTemplates) > 0 {
			def += "\tswitch {\n"
			for expr := range msg.CustomTemplates {
				def += fmt.Sprintf("\tcase %s:\n", expr)
				// will index proper lang file later
				def += fmt.Sprintf("\t\tdata.I18ngoCustomTemplateExpr = \"%s\"\n", expr)
				def += "\t\treturn t.l.RenderMessage(\"" + msgID + "\", data)\n"
			}
			def += "\t}\n"
		}

		def += fmt.Sprintf("\treturn t.l.RenderMessage(\"%s\", data)\n}\n", msgID)

		buf.WriteString(def)
	}

	source := buf.Bytes()

	source, err = format.Source(source, format.Options{})
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

// RenderMessage renders a message based on the current language and id.
func (l *LanguageLoader) RenderMessage(ctx context.Context, id string, data interface{}) (string, error) {
	lang, _ := ctx.Value(LangCtxKey{}).(string) // TODO: lang should just be set on New() and reused for all messages.
	translations, ok := l.translations[lang]
	if !ok {
		return "", errors.New("language not loaded")
	}

	msg, exists := translations.Messages[id]
	if !exists {
		return "", errors.New("message id not found")
	}

	// Render the message template with the provided data
	tmpl, err := template.New("message").Parse(msg.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
