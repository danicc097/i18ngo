package main

import (
	"context"
	"os"

	"github.com/danicc097/i18ngo/templates"
)

func main() {
	component := templates.TranslationCode(templates.TemplateData{
		Langs: []templates.LangData{
			{Lang: "en"},
		},
	})
	component.Render(context.Background(), os.Stdout)
}
