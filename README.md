# i18ngo

**WIP**

Generate simple i18n go code based on a schema-backed `{lang}.i18ngo.yaml`:

```yaml
messages:
  my_greeting:
    template: "Hello {{ .Name }}! You have {{ .Count }} messages."
    variables:
      Name: string
      Count: int
    custom_templates:
      # any valid go expression is allowed. Vars are available in camelCase form.
      - expression: "count == 1"
        template: "Hello {{ .Name }}! You have {{ .Count }} message."
      - expression: "count == 0"
        template: "Hello {{ .Name }}! You have no messages."
```

The above will generate code you can use in html templates with
 `*.MyGreeting(count, name)` with the current loader. Using alongside a library
 like `a-h/templ`,
 messages also benefit from full LSP support.

Initialize all translators at startup from the generated code:

```go
// call as library to generate as many packages as you want.
// include a templates/template.go.tpl in your filesystem, starting
// with the official one as basepoint and extend if required.
err := i18n.Generate(fsys, "path to *.i18ngo.yaml dir", pkgName)

// assuming your codegen was saved to an i18ngen package
tt := i18ngen.NewTranslators()

// lang may come from context, etc.
t := tt[lang]
t.MyGreeting(count, name)
```
