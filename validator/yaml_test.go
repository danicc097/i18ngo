package validator_test

import (
	"testing"
	"testing/fstest"

	"github.com/danicc097/i18ngo/validator"
)

func TestCompareTranslationFiles(t *testing.T) {
	testCases := []struct {
		name          string
		files         map[string]string
		expectedError bool
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
			expectedError: false,
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
			expectedError: true,
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
			if (err != nil) != tc.expectedError {
				t.Errorf("CompareTranslationFiles() error = %v, expectedError %v", err, tc.expectedError)
			}
		})
	}
}
