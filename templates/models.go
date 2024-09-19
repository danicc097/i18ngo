package templates

type TemplateData struct {
	PkgName      string
	Langs        []LangData
	Messages     []MessageData
	Translations []TranslationData
}

type LangData struct {
	CamelLang string
	Lang      string
}

type MessageData struct {
	CamelLang       string
	MethodName      string
	Args            string
	Vars            []VarData
	Template        string
	CustomTemplates []CustomTemplate
}

type VarData struct {
	Name  string
	Type  string
	Param string
}

type TranslationData struct {
	CamelLang string
	Messages  []MessageData
}

type CustomTemplate struct {
	Expression string `yaml:"expression"`
	Template   string `yaml:"template"`
}

type Message struct {
	Template        string            `yaml:"template"`
	Variables       map[string]string `yaml:"variables"`
	CustomTemplates []CustomTemplate  `yaml:"custom_templates"`
}

type Translations struct {
	Messages map[string]Message `yaml:"messages"`
}
