package translation

import (
	"bytes"
	"html/template"

	"github.com/danicc097/i18ngo"
)

// Translator is implemented by all language translators.
type Translator interface {
	MyGreeting(age int, name string) (string, error)
}

// Lang represents available translated languages.
type Lang string

const (
	LangEn Lang = "en"
	LangEs Lang = "es"
)

// NewTranslators initializes all translators based on a initialized loader.
func NewTranslators(l *i18ngo.LanguageLoader) map[Lang]Translator {
	return map[Lang]Translator{
		LangEn: newEn(l),
		LangEs: newEs(l),
	}
}

type En struct {
	l *i18ngo.LanguageLoader
}

func newEn(l *i18ngo.LanguageLoader) *En {
	return &En{
		l: l,
	}
}

// MyGreeting renders a properly translated message.
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

func newEs(l *i18ngo.LanguageLoader) *Es {
	return &Es{
		l: l,
	}
}

// MyGreeting renders a properly translated message.
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
