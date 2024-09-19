package i18ngo_test

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

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
		got, err := i18ngo.Generate(testValidFS, testName, pkgName)
		if err != nil {
			t.Fatalf("Failed to generate Go code for %s/: %v", entry.Name(), err)
		}

		wantSnapshot := filepath.Join(testName, "snapshots", "i18n.go")
		want, err := os.ReadFile(wantSnapshot)
		if err != nil {
			t.Fatalf("Failed to read snapshot file for %s: %v", entry.Name(), err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch in %q (-want +got):\n%s", testdataDir+"/"+entry.Name(), diff)
			t.Log(string(got))
		}
	}
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

		wantErrorPath := filepath.Join(testName, "stderr.txt")
		e, err := testInvalidFS.ReadFile(wantErrorPath)
		require.NoError(t, err)
		wantError := strings.TrimSuffix(string(e), "\n")

		_, err = i18ngo.Generate(testInvalidFS, testName, pkgName)
		if err == nil {
			t.Fatalf("Expected error for %s/ but got nothing", entry.Name())
		}

		assert.ErrorContainsf(t, err, string(wantError), entry.Name())
	}
}

func TestWithCustomTemplate(t *testing.T) {
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

	got, err := i18ngo.Generate(fs, testdataDir, pkgName, i18ngo.WithFilesystemTemplate())
	require.NoError(t, err)
	assert.Equal(t, "package customtemplate\n", string(got))
}
