package translation

import (
	"bytes"
	"html/template"
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

// NewTranslators initializes all translators.
func NewTranslators() map[Lang]Translator {
	return map[Lang]Translator{
		LangEn: newEn(),
		LangEs: newEs(),
	}
}

type en struct{}

func newEn() *en {
	return &en{}
}

// MyGreeting renders a properly translated message.
func (t *en) MyGreeting(age int, name string) (string, error) {
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

type es struct{}

func newEs() *es {
	return &es{}
}

// MyGreeting renders a properly translated message.
func (t *es) MyGreeting(age int, name string) (string, error) {
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
