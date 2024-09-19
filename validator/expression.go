package validator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
)

func ValidateCustomExpression(expression string, variables []string) error {
	fs := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fs, "", expression, 0)
	if err != nil {
		return fmt.Errorf("invalid go expression: %w", err)
	}

	if len(variables) == 0 {
		return nil
	}

	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if !slices.Contains(variables, ident.Name) {
				err = fmt.Errorf("unknown variable used in expression: %s", ident.Name)
				return false
			}
		}
		return true
	})

	return err
}
