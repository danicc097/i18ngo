# i18ngo

**WIP**

Generate simple i18n go code based on a schema-backed `{lang}.i18ngo.yaml`:

```yaml
messages:
  my_greeting:
    template: "Hello {{ .Name }}! You have {{ .Count }} messages."
    variables:
      # vars are ensured to all be defined upon codegen
      Name: string
      Count: int
    custom_templates:
      # any valid go expression is allowed. Vars are available in camelCase form.
      "count == 0":  "Hello {{ .Name }}! You have {{ .Count }} message."
```

The above will generate code you can use in html templates as `*.MyGreeting(age,
name)` with the current loader. Any other language
