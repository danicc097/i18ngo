package validator_test

import (
	"testing"
	"testing/fstest"

	"github.com/danicc097/i18ngo/validator"
	"github.com/stretchr/testify/require"
)

func TestCompareTranslationFiles(t *testing.T) {
	testCases := []struct {
		name      string
		files     map[string]string
		wantError string
	}{
		{
			name: "Matching structures",
			files: map[string]string{
				"data/en.i18ngo.yaml": `messages:
  my_greeting:
    template: "a"
    variables:
      Name: string
      Count: int
    custom_templates:
      "count == 0":  "a"`,
				"data/es.i18ngo.yaml": `messages:
  my_greeting:
    template: "b"
    variables:
      Name: string
      Count: int
    custom_templates:
      "count == 0": "b"`,
			},
		},
		{
			name: "Mismatched structures",
			files: map[string]string{
				"data/en.i18ngo.yaml": `messages:
  my_greeting:
    template: "a"
    variables:
      Name: string
      Count: int
    custom_templates:
      "count == 0":  "a"`,
				"data/es.i18ngo.yaml": `messages:
  my_greeting:
    template: "b"
    variables:
      Name: string
      Count: int
    custom_templates:
      "count == 10000": "b"`,
			},
			wantError: `structure mismatch between translation files "data/en.i18ngo.yaml" and "data/es.i18ngo.yaml" at .messages.my_greeting.custom_templates.count == 0`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fsys := fstest.MapFS{}
			for name, content := range tc.files {
				fsys[name] = &fstest.MapFile{
					Data: []byte(content),
				}
			}

			err := validator.ValidateTranslationFiles(fsys, "data")
			if tc.wantError != "" {
				require.EqualError(t, err, tc.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
