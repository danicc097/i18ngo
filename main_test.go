package i18ngo_test

import (
	"embed"
	"go/format"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	custom_template_t "github.com/danicc097/i18ngo/testdata/valid/custom_template/snapshots"

	"github.com/danicc097/i18ngo"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/valid/*
//go:embed templates/template.go.tpl
var testValidFS embed.FS

//go:embed testdata/invalid/*
//go:embed templates/template.go.tpl
var testInvalidFS embed.FS

const pkgName = "translations"

func TestCodeGeneration(t *testing.T) {
	t.Parallel()

	testdataDir := "testdata/valid"
	entries, err := testValidFS.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testName := filepath.Join(testdataDir, entry.Name())
		data, err := i18ngo.GetTranslationData(testValidFS, testName, pkgName)
		require.NoError(t, err)
		got, err := i18ngo.Generate(data)
		if err != nil {
			t.Fatalf("Failed to generate Go code for %s/: %v", entry.Name(), err)
		}

		wantSnapshot := filepath.Join(testName, "snapshots", "i18n.go")
		// format both src with gofmt:
		want, err := os.ReadFile(wantSnapshot) // don't use fsys for tests, snapshot will be updated later
		if err != nil {
			t.Fatalf("Failed to read snapshot file for %s: %v", entry.Name(), err)
		}
		wantFmtted := mustFormat(t, want)
		gotFmtted := mustFormat(t, got)
		if diff := cmp.Diff(string(wantFmtted), string(gotFmtted)); diff != "" {
			t.Errorf("Mismatch in %q (-want +got):\n%s", testdataDir+"/"+entry.Name(), diff)
		}

		if os.Getenv("SNAPSHOT_UPDATE") != "" {
			if err := os.WriteFile(wantSnapshot, got, 0o666); err != nil {
				t.Fatalf("Failed to update snapshot file for %s: %v", entry.Name(), err)
			}
		}
	}
}

func mustFormat(t *testing.T, src []byte) []byte {
	t.Helper()

	got, err := format.Source(src)
	if err != nil {
		t.Fatalf("Failed to format generated code: %v", err)
	}
	return got
}

func TestInvalidCodeGeneration(t *testing.T) {
	t.Parallel()

	testdataDir := "testdata/invalid"
	entries, err := testInvalidFS.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testName := filepath.Join(testdataDir, entry.Name())

		t.Run(testName, func(t *testing.T) {
			wantErrorPath := filepath.Join(testName, "stderr.txt")
			e, err := testInvalidFS.ReadFile(wantErrorPath)
			require.NoError(t, err)
			wantError := strings.TrimSuffix(string(e), "\n")
			data, err := i18ngo.GetTranslationData(testInvalidFS, testName, pkgName)
			if err == nil {
				_, err = i18ngo.Generate(data)
				if err == nil {
					t.Fatalf("Expected error for %s/ but got nothing", entry.Name())
				}
			}

			assert.ErrorContainsf(t, err, string(wantError), entry.Name())
		})
	}
}

func TestWithCustomTemplate(t *testing.T) {
	t.Skip(`
	TODO: allow to call custom templ generated code (provide a Generator interface).
	This way users can use whatever templ version they want, or plain go code for all we care.
	`)

	t.Parallel()

	testdataDir := "testdata"
	fs := fstest.MapFS{
		"templates/template.go.tpl": &fstest.MapFile{
			Data: []byte("package customtemplate"),
		},
		"testdata/en.i18ngo.yaml": &fstest.MapFile{
			Data: []byte(`messages:
  my_greeting:
    template: "a"
    variables:
      Name: string
      Count: int`),
		},
	}

	data, err := i18ngo.GetTranslationData(fs, testdataDir, pkgName, i18ngo.WithFilesystemTemplate())
	require.NoError(t, err)

	got, err := i18ngo.Generate(data)
	require.NoError(t, err)
	assert.Equal(t, "package customtemplate\n", string(got))
}

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

	for _, tc := range testCases {
		t.Run(string(tc.lang)+"_"+tc.name+"_"+strconv.Itoa(tc.count), func(t *testing.T) {
			tr := tt[tc.lang]
			tr = custom_template_t.NewMemoizedTranslator(tr)
			out, err := tr.MyGreeting(tc.count, tc.name)
			require.NoError(t, err)
			require.Equal(t, tc.expected, out)
		})
	}
}
