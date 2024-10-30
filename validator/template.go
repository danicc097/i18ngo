package validator

import (
	"fmt"
	"html/template"
	"regexp"
	"slices"
	"strings"

	"github.com/kenshaw/snaker"
)

// ValidateTemplate checks if all variables used inside {{ .MyVar }} exist in the provided variables.
func ValidateTemplate(tpl string, variables []string) error {
	nvars := make([]string, 0, len(variables))
	for _, v := range variables {
		nvars = append(nvars, snaker.ForceLowerCamelIdentifier(v))
	}
	errors := []string{}

	if _, err := template.New("").Parse(tpl); err != nil {
		return fmt.Errorf("unparseable template: %w", err)
	}

	validRe := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)

	generalRe := regexp.MustCompile(`\{\{[^}]*\}\}`)
	generalMatches := generalRe.FindAllString(tpl, -1)

	errorMatches := regexp.MustCompile(`\{\s*[.]?[^\s][^}]*\}\}`).FindAllString(tpl, -1)
	for _, errMatch := range errorMatches {
		if strings.HasPrefix(errMatch, "{{") { // poor mans neg lookahead
			continue
		}
		errors = append(errors, fmt.Sprintf("possible invalid syntax: %s", errMatch))
	}
	for _, match := range generalMatches {
		matches := validRe.FindAllStringSubmatch(match, -1)
		if len(matches) > 0 {
			varName := matches[0][1]
			if !slices.Contains(variables, varName) {
				errors = append(errors, fmt.Sprintf("unknown variable used in template: %s", varName))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("invalid template: %s", strings.Join(errors, ", "))
	}

	return nil
}
