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

func (t *T) MyGreeting(count int, name string) (string, error) {
	data := struct {
		Count                    int
		Name                     string
		I18ngoCustomTemplateExpr string
	}{
		Count: count,
		Name:  name,
	}
	switch {
	case count == 1:
		data.I18ngoCustomTemplateExpr = "count == 1"
		return t.l.RenderMessage("my_greeting", data)
	case count == 0:
		data.I18ngoCustomTemplateExpr = "count == 0"
		return t.l.RenderMessage("my_greeting", data)
	}
	return t.l.RenderMessage("my_greeting", data)
}
