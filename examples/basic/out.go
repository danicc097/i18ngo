// Code generated by i18ngo. DO NOT EDIT.
package examples_basic

import (
	"bytes"
	"html/template"
)

// Translator is implemented by all language translators.
type Translator interface {
	MyGreeting(count int, name string) (string, error)
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
func (t *en) MyGreeting(count int, name string) (string, error) {
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

type es struct{}

func newEs() *es {
	return &es{}
}

// MyGreeting renders a properly translated message.
func (t *es) MyGreeting(count int, name string) (string, error) {
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
