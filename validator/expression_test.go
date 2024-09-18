package validator_test

import (
	"testing"

	"github.com/danicc097/i18ngo/validator"
)

// NOTE: already extensively tested in stdlib.
func TestValidateGoExpression(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		shouldErr bool
	}{
		{"ValidExpression1", "count == 0", false},
		{"ValidExpression2", "name != \"\"", false},
		{"InvalidExpression", "invalid-expression(", true},
		{"ValidExpression3", "age > 18 && count != 0", false},
		{"InvalidExpression2", "count !@ 0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateGoExpression(tt.expr)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateGoExpression() error = %v, wantErr %v", err, tt.shouldErr)
			}
		})
	}
}
