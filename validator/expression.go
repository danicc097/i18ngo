package validator

import (
	"go/parser"
	"go/token"
)

func ValidateGoExpression(expr string) error {
	fs := token.NewFileSet()
	_, err := parser.ParseExprFrom(fs, "", expr, 0)
	return err
}
