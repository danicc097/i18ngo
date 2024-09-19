package translation

import (
	"bytes"
	"html/template"

	"github.com/danicc097/i18ngo"
)

type Translator interface {
	MyGreeting(age int, name string) (string, error)
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

func (t *En) MyGreeting(age int, name string) (string, error) {
	data := struct {
		Age  int
		Name string
	}{
		Age:  age,
		Name: name,
	}
	tmpl, err := template.New("message").Parse("Hello {{ .Name }}! You are {{ .Age }} years old.")
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

func (t *Es) MyGreeting(age int, name string) (string, error) {
	data := struct {
		Age  int
		Name string
	}{
		Age:  age,
		Name: name,
	}
	tmpl, err := template.New("message").Parse("Hola {{ .Name }}! Tienes {{ .Age }} años.")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
