package translation

import (
	"bytes"
	"html/template"
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

// NewTranslators initializes all translators.
func NewTranslators() map[Lang]Translator {
	return map[Lang]Translator{
{{- range .Langs }}
		Lang{{.CamelLang}}: new{{.CamelLang}}(),
{{- end }}
	}
}

{{- range .Translations }}
type {{camelCase .CamelLang}} struct {}

func new{{.CamelLang}}() *{{camelCase .CamelLang}} {
	return &{{camelCase .CamelLang}}{}
}

{{- range .Messages }}
// {{.MethodName}} renders a properly translated message.
func (t *{{camelCase .CamelLang}}) {{.MethodName}}({{.Args}}) (string, error) {
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
