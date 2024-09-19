package translation

import (
	"bytes"
	"html/template"

	"github.com/danicc097/i18ngo"
)

type Translator interface {
	MyGreeting(count int, name string) (string, error)
}

type Lang string

const (
	LangEn Lang = "en"
	LangEs Lang = "es"
)

var Translators = map[Lang]Translator{
	LangEn: NewEn(&i18ngo.LanguageLoader{}),
	LangEs: NewEs(&i18ngo.LanguageLoader{}),
}

type En struct {
	l *i18ngo.LanguageLoader
}

func NewEn(l *i18ngo.LanguageLoader) *En {
	return &En{
		l: l,
	}
}

func (t *En) MyGreeting(count int, name string) (string, error) {
	data := struct {
		Count int
		Name  string
	}{
		Count: count,
		Name:  name,
	}
	switch {
	case count == 1:
		tmpl, err := template.New("custom").Parse("Hello {{ .Name }}! You have {{ .Count }} message.")
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	case count == 0:
		tmpl, err := template.New("custom").Parse("Hello {{ .Name }}! You have no messages.")
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	tmpl, err := template.New("message").Parse("Hello {{ .Name }}! You have {{ .Count }} messages.")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type Es struct {
	l *i18ngo.LanguageLoader
}

func NewEs(l *i18ngo.LanguageLoader) *Es {
	return &Es{
		l: l,
	}
}

func (t *Es) MyGreeting(count int, name string) (string, error) {
	data := struct {
		Count int
		Name  string
	}{
		Count: count,
		Name:  name,
	}
	switch {
	case count == 1:
		tmpl, err := template.New("custom").Parse("Hola {{ .Name }}! Tienes {{ .Count }} mensaje.")
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	case count == 0:
		tmpl, err := template.New("custom").Parse("Hola {{ .Name }}! No tienes ningún mensaje.")
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	tmpl, err := template.New("message").Parse("Hola {{ .Name }}! Tienes {{ .Count }} mensajes.")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
