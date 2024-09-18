package i18ngo_test

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/danicc097/i18ngo"
	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/valid/*
var testValidFS embed.FS

func TestCodeGeneration(t *testing.T) {
	testdataDir := "testdata/valid"
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		langDir := filepath.Join(testdataDir, entry.Name())
		code, err := i18ngo.Generate(testValidFS, entry.Name())
		if err != nil {
			t.Fatalf("Failed to generate Go code for %s: %v", entry.Name(), err)
		}

		wantSnapshot := filepath.Join(langDir, "snapshots", "i18n.go")
		wantBytes, err := os.ReadFile(wantSnapshot)
		if err != nil {
			t.Fatalf("Failed to read snapshot file for %s: %v", entry.Name(), err)
		}

		if diff := cmp.Diff(wantBytes, code); diff != "" {
			t.Errorf("Mismatch between generated code and snapshot for %s (-want +got):\n%s", entry.Name(), diff)
		}
	}
}