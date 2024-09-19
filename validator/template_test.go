package validator_test

import (
	"testing"

	"github.com/danicc097/i18ngo/validator"
	"github.com/stretchr/testify/require"
)

func TestValidateTemplate(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		variables []string
		wantError string
	}{
		{
			name:      "valid template with matching variables",
			template:  `{{ .MyVar }} and {{ .AnotherVar }}`,
			variables: []string{"MyVar", "AnotherVar"},
		},
		{
			name:      "invalid template with unknown variable",
			template:  `{{ .MyVar }} and {{ .UnknownVar }}`,
			variables: []string{"MyVar", "AnotherVar"},
			wantError: "unknown variable used in template: UnknownVar",
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
			name:      "template with invalid variable format",
			template:  `{{ .123InvalidVar }}`,
			variables: []string{"MyVar"},
			wantError: "invalid template: template: :1: bad number syntax: \".123I\"",
		},
		{
			name:      "template with multiple unknown variables",
			template:  `{{ .UnknownVar1 }} and {{ .UnknownVar2 }}`,
			variables: []string{"MyVar"},
			wantError: "unknown variable used in template: UnknownVar1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTemplate(tt.template, tt.variables)
			if tt.wantError != "" {
				require.EqualError(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
