package validator_test

import (
	"strconv"
	"testing"

	custom_template_t "github.com/danicc097/i18ngo/testdata/valid/custom_template/snapshots"
	"github.com/stretchr/testify/require"
)

func TestTranslationsCustomTemplate(t *testing.T) {
	t.Parallel()

	tt := custom_template_t.NewTranslators()

	type testCase struct {
		lang     custom_template_t.Lang
		count    int
		name     string
		expected string
	}

	testCases := []testCase{
		{
			lang:     custom_template_t.LangEn,
			count:    0,
			name:     "Alice",
			expected: "Hello Alice! You have no messages.",
		},
		{
			lang:     custom_template_t.LangEn,
			count:    1,
			name:     "Bob",
			expected: "Hello Bob! You have 1 message.",
		},
		{
			lang:     custom_template_t.LangEn,
			count:    2,
			name:     "Charlie",
			expected: "Hello Charlie! You have 2 messages.",
		},
		{
			lang:     custom_template_t.LangEs,
			count:    0,
			name:     "Ana",
			expected: "Hola Ana! No tienes ningún mensaje.",
		},
		{
			lang:     custom_template_t.LangEs,
			count:    1,
			name:     "Juan",
			expected: "Hola Juan! Tienes 1 mensaje.",
		},
		{
			lang:     custom_template_t.LangEs,
			count:    10,
			name:     "Luis",
			expected: "Hola Luis! Tienes 10 mensajes.",
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.lang)+"_"+tc.name+"_"+strconv.Itoa(tc.count), func(t *testing.T) {
			tr := tt[tc.lang]
			out, err := tr.MyGreeting(tc.count, tc.name)
			require.NoError(t, err)
			require.Equal(t, tc.expected, out)
		})
	}
}
