package validator_test

import (
	"testing"

	"github.com/danicc097/i18ngo/validator"
	"github.com/stretchr/testify/require"
)

func TestValidateTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		variables   []string
		errContains string
	}{
		{
			name:      "valid template with matching variables",
			template:  `{{ .MyVar }} and {{ .AnotherVar }}`,
			variables: []string{"MyVar", "AnotherVar"},
		},
		{
			name:        "template.Parse edge case not considered syntax error",
			template:    `{ .MyVar }} and {{ .MyVar }}`,
			variables:   []string{"MyVar", "AnotherVar"},
			errContains: "invalid template: possible invalid syntax: { .MyVar }}",
		},
		{
			name:        "template.Parse 1",
			template:    `{{ .MyVar } and {{ .MyVar }}`,
			variables:   []string{"MyVar", "AnotherVar"},
			errContains: "unparseable template",
		},
		{
			name:        "template.Parse 2",
			template:    `{{ .MyVar and {{ .MyVar }}`,
			variables:   []string{"MyVar", "AnotherVar"},
			errContains: "unparseable template",
		},
		{
			name:        "template.Parse 3",
			template:    `{{ .123InvalidVar }}`,
			variables:   []string{"MyVar"},
			errContains: "unparseable template",
		},
		{
			name:        "invalid template with unknown variable",
			template:    `{{ .MyVar }} and {{ .UnknownVar }}`,
			variables:   []string{"MyVar", "AnotherVar"},
			errContains: "unknown variable used in template: UnknownVar",
		},
		{
			name:      "template with no variables",
			template:  `No variables here`,
			variables: []string{"MyVar", "AnotherVar"},
		},
		{
			name:      "valid template with extra variables",
			template:  `{{ .MyVar }}`,
			variables: []string{"MyVar", "ExtraVar"},
		},
		{
			name:        "template with multiple unknown variables",
			template:    `{{ .UnknownVar1 }} and {{ .UnknownVar2 }}`,
			variables:   []string{"MyVar"},
			errContains: "unknown variable used in template: UnknownVar1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTemplate(tt.template, tt.variables)
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
