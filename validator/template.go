package validator

import (
	"fmt"
	"html/template"
	"regexp"
	"slices"
)

// ValidateTemplate checks if all variables used inside {{ .MyVar }} exist in the provided variables.
func ValidateTemplate(tpl string, variables []string) error {
	if _, err := template.New("").Parse(tpl); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	re := regexp.MustCompile(`\{\{\s*\.([a-zA-Z0-9_]*)\s*\}\}`)
	matches := re.FindAllStringSubmatch(tpl, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue // no capture group
		}
		varName := match[1]

		if !slices.Contains(variables, varName) {
			return fmt.Errorf("unknown variable used in template: %s", varName)
		}
	}

	return nil
}
