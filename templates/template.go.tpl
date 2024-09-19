package translation

import (
	"bytes"
	"html/template"

	"github.com/danicc097/i18ngo"
)

// Translator is implemented by all language translators.
type Translator interface {
{{- range .Messages }}
	{{.MethodName}}({{.Args}}) (string, error)
{{- end }}
}

// Lang represents available translated languages.
type Lang string

const (
{{- range .Langs }}
	Lang{{.CamelLang}} Lang = "{{.Lang}}"
{{- end }}
)

// NewTranslators initializes all translators based on a initialized loader.
func NewTranslators(l *i18ngo.LanguageLoader) map[Lang]Translator {
	return map[Lang]Translator{
{{- range .Langs }}
		Lang{{.CamelLang}}: new{{.CamelLang}}(l),
{{- end }}
	}
}

{{- range .Translations }}
type {{.CamelLang}} struct {
	l *i18ngo.LanguageLoader
}

func new{{.CamelLang}}(l *i18ngo.LanguageLoader) *{{.CamelLang}} {
	return &{{.CamelLang}}{
		l: l,
	}
}

{{- range .Messages }}
// {{.MethodName}} renders a properly translated message.
func (t *{{.CamelLang}}) {{.MethodName}}({{.Args}}) (string, error) {
	data := struct {
	{{- range .Vars }}
		{{.Name}} {{.Type}}
	{{- end }}
	}{
	{{- range .Vars }}
		{{.Name}}: {{.Param}},
	{{- end }}
	}

{{- if .CustomTemplates }}
	switch {
	{{- range .CustomTemplates }}
	case {{.Expression}}:
		tmpl, err := template.New("custom").Parse("{{.Template}}")
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	{{- end }}
	}
{{- end }}

	tmpl, err := template.New("message").Parse("{{.Template}}")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
{{- end }}
{{- end }}
