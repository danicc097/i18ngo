package validator_test

import (
	"testing"

	"github.com/danicc097/i18ngo/validator"
)

func TestValidateCustomTemplate(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		variables []string
		shouldErr bool
	}{
		{"ValidExpression1", "count == 0", []string{"count"}, false},
		{"ValidExpression2", "name != \"\"", []string{"name"}, false},
		{"ValidExpression3", "age > 18 && count != 0", []string{"age", "count"}, false},

		{"InvalidVariable1", "unknown == 0", []string{"count"}, true},
		{"InvalidVariable2", "name != \"\" && id > 0", []string{"name"}, true},

		// Syntax errors (already extensively tested in std lib)
		{"InvalidExpression", "invalid-expression(", []string{}, true},
		{"InvalidSyntax2", "count !@ 0", []string{"count"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCustomTemplate(tt.expr, tt.variables)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateCustomTemplate() error = %v, wantErr %v", err, tt.shouldErr)
			}
		})
	}
}
