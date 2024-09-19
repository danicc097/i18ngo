package i18ngo

import (
	"bytes"
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

type CustomTemplate struct {
	Expression string `yaml:"expression"`
	Template   string `yaml:"template"`
}

type Message struct {
	Template        string            `yaml:"template"`
	Variables       map[string]string `yaml:"variables"`
	CustomTemplates []CustomTemplate  `yaml:"custom_templates"`
}

type Translations struct {
	Messages map[string]Message `yaml:"messages"`
}

type LanguageLoader struct {
	translations map[string]Translations
}

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

func Generate(fsys fs.FS, path string) ([]byte, error) {
	loader, err := NewLanguageLoader(fsys, path)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	buf.WriteString("package translation\n\n")
	buf.WriteString("import (\n\t\"bytes\"\n\t\"html/template\"\n\t\"github.com/danicc097/i18ngo\"\n)\n\n")

	buf.WriteString("type Translator interface {\n")

	firstLang := ""
	for k := range loader.translations {
		firstLang = k
		break
	}

	messageIDs := make([]string, 0, len(loader.translations[firstLang].Messages))
	for id := range loader.translations[firstLang].Messages {
		messageIDs = append(messageIDs, id)
	}
	sort.Strings(messageIDs)

	for _, id := range messageIDs {
		methName := snaker.SnakeToCamel(id)
		vars := loader.translations[firstLang].Messages[id].Variables
		varsList := make([]string, 0, len(vars))
		for name := range vars {
			varsList = append(varsList, name)
		}
		sort.Strings(varsList)
		varArgs := ""
		for _, name := range varsList {
			typ := vars[name]
			varArgs += fmt.Sprintf("%s %s, ", snaker.ForceLowerCamelIdentifier(name), typ)
		}
		varArgs = strings.TrimSuffix(varArgs, ", ")
		buf.WriteString(fmt.Sprintf("\t%s(%s) (string, error)\n", methName, varArgs))
	}
	buf.WriteString("}\n\n")

	langs := make([]string, 0, len(loader.translations))
	for lang := range loader.translations {
		langs = append(langs, lang)
	}
	sort.Strings(langs)

	buf.WriteString("type Lang string\n")
	buf.WriteString("const (\n")
	for _, lang := range langs {
		camelLang := snaker.SnakeToCamel(lang)
		buf.WriteString(fmt.Sprintf("Lang%s Lang = \"%s\"\n", camelLang, lang))
	}
	buf.WriteString(")\n\n")

	buf.WriteString("var Translators = map[Lang]Translator{\n")
	for _, lang := range langs {
		camelLang := snaker.SnakeToCamel(lang)
		buf.WriteString(fmt.Sprintf("\tLang%s: New%s(&i18ngo.LanguageLoader{}),\n", camelLang, camelLang))
	}
	buf.WriteString("}\n\n")

	tt := make([]string, 0, len(loader.translations))
	for k := range loader.translations {
		tt = append(tt, k)
	}
	sort.Strings(tt)

	for _, lang := range tt {
		translations := loader.translations[lang]
		camelLang := snaker.SnakeToCamel(lang)
		buf.WriteString(fmt.Sprintf("type %s struct {\n\tl *i18ngo.LanguageLoader\n}\n\n", camelLang))
		buf.WriteString(fmt.Sprintf("func New%s(l *i18ngo.LanguageLoader) *%s {\n\treturn &%s{\n\t\tl: l,\n\t}\n}\n\n", camelLang, camelLang, camelLang))

		messageIDs := make([]string, 0, len(translations.Messages))
		for id := range translations.Messages {
			messageIDs = append(messageIDs, id)
		}
		sort.Strings(messageIDs)

		for _, msgID := range messageIDs {
			msg := translations.Messages[msgID]
			methName := snaker.SnakeToCamel(msgID)

			def := fmt.Sprintf("func (t *%s) %s(", camelLang, methName)

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
			def += "\t} {\n"
			for _, name := range varNames {
				def += fmt.Sprintf("\t\t%s: %s,\n", snaker.SnakeToCamel(name), snaker.ForceLowerCamelIdentifier(name))
			}
			def += "\t}\n"

			if len(msg.CustomTemplates) > 0 {
				def += "\tswitch {\n"
				for _, ct := range msg.CustomTemplates {
					def += fmt.Sprintf("\tcase %s:\n", ct.Expression)
					def += fmt.Sprintf("\t\ttmpl, err := template.New(\"custom\").Parse(\"%s\")\n", template.HTMLEscapeString(ct.Template))
					def += "\t\tif err != nil {\n\t\t\treturn \"\", err\n\t\t}\n"
					def += "\t\tvar buf bytes.Buffer\n"
					def += "\t\tif err := tmpl.Execute(&buf, data); err != nil {\n\t\t\treturn \"\", err\n\t\t}\n"
					def += "\t\treturn buf.String(), nil\n"
				}
				def += "\t}\n"
			}

			def += "\ttmpl, err := template.New(\"message\").Parse(\"" + template.HTMLEscapeString(msg.Template) + "\")\n"
			def += "\tif err != nil {\n\t\treturn \"\", err\n\t}\n"
			def += "\tvar buf bytes.Buffer\n"
			def += "\tif err := tmpl.Execute(&buf, data); err != nil {\n\t\treturn \"\", err\n\t}\n"
			def += "\treturn buf.String(), nil\n}\n"

			buf.WriteString(def)
		}
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
