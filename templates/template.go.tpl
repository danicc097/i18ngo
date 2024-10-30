{{- /* gotype: github.com/danicc097/i18ngo/templates.TemplateData */ -}}
// Code generated by i18ngo. DO NOT EDIT.
package {{ .PkgName }}

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
type {{camelCase .CamelLang}} struct {
    {{- range .Messages }}
    {{ .MethodName }}Dft *template.Template
    {{- if .CustomTemplates }}
        {{- $methodName := .MethodName }}
        {{- range $index, $ct := .CustomTemplates }}
    {{ $methodName }}Custom{{ $index }} *template.Template
        {{- end }}
    {{- end }}
    {{- end }}
}

func new{{.CamelLang}}() *{{camelCase .CamelLang}} {
    return &{{camelCase .CamelLang}}{
    {{- range .Messages }}
        {{ .MethodName }}Dft: template.Must(template.New("{{ .MethodName }}").Parse("{{ .Template }}")),
        {{- if .CustomTemplates }}
            {{- $methodName := .MethodName }}
            {{- range $index, $ct := .CustomTemplates }}
        {{ $methodName }}Custom{{ $index }}: template.Must(template.New("{{ $methodName }}Custom{{ $index }}").Parse("{{ $ct.Template }}")),
            {{- end }}
        {{- end }}
    {{- end }}
    }
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
    var tmpl *template.Template
    {{- if .CustomTemplates }}
    switch {
        {{- $methodName := .MethodName }}
        {{- range $index, $ct := .CustomTemplates }}
    case {{ $ct.Expression }}:
        tmpl = t.{{ $methodName }}Custom{{ $index }}
        {{- end }}
    default:
        tmpl = t.{{ .MethodName }}Dft
    }
    {{- else }}
    tmpl = t.{{ .MethodName }}Dft
    {{- end }}
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}
{{- end }}
{{- end }}
