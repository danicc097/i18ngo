package translation

import (
	"context"

	"github.com/danicc097/i18ngo"
)

type T struct {
	Ctx context.Context
	l   *i18ngo.LanguageLoader
}

// New returns a new i18n translator.
func New(l *i18ngo.LanguageLoader) *T {
	return &T{
		l: l,
	}
}

func (t *T) MyGreeting(age int, name string) (string, error) {
	data := struct {
		Age                      int
		Name                     string
		I18ngoCustomTemplateExpr string
	}{
		Age:  age,
		Name: name,
	}
	return t.l.RenderMessage("my_greeting", data)
}
