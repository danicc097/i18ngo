package validator

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

type anyMap map[string]interface{}

// ValidateTranslationFiles verifies the structure of translation files in the given path is the same.
func ValidateTranslationFiles(fsys fs.FS, path string) error {
	// NOTE: variable type leaf nodes not checked since interface wont be implemented by that language codegen anyway.

	var files []string
	err := fs.WalkDir(fsys, path, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(name, ".i18ngo.yaml") {
			files = append(files, name)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	var structures []anyMap
	for _, file := range files {
		structure, err := parseYAMLFile(fsys, file)
		if err != nil {
			return fmt.Errorf("error parsing YAML file %q: %w", file, err)
		}
		structures = append(structures, structure)
	}

	// Compare each file structure with the first one
	for i := 1; i < len(structures); i++ {
		if ok, diffPath := compareMaps(structures[0], structures[i], ""); !ok {
			return fmt.Errorf("structure mismatch between translation files %q and %q at %s", files[0], files[i], diffPath)
		}
	}

	return nil
}

func parseYAMLFile(fsys fs.FS, filename string) (anyMap, error) {
	file, err := fsys.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content anyMap
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&content); err != nil {
		return nil, err
	}
	return content, nil
}

// compareMaps compares two maps recursively, returning false if there is a difference.
// It also returns the path where the difference occurs.
func compareMaps(map1, map2 anyMap, currentPath string) (bool, string) {
	if len(map1) != len(map2) {
		return false, currentPath
	}

	for key1, val1 := range map1 {
		val2, exists := map2[key1]
		if !exists {
			return false, fmt.Sprintf("%s.%s", currentPath, key1)
		}

		if ok, diffPath := compareValues(val1, val2, fmt.Sprintf("%s.%s", currentPath, key1)); !ok {
			return false, diffPath
		}
	}

	return true, ""
}

// compareValues compares two values, considering map or slice types, returning false if they differ.
func compareValues(val1, val2 interface{}, currentPath string) (bool, string) {
	map1, ok1 := val1.(anyMap)
	map2, ok2 := val2.(anyMap)
	if ok1 && ok2 {
		return compareMaps(map1, map2, currentPath)
	}

	// If both are slices, we don't check contents, just ensure both are slices
	_, ok1 = val1.([]interface{})
	_, ok2 = val2.([]interface{})
	if ok1 && ok2 {
		return true, "" // Slices are not compared, just their type
	}

	if ok1 != ok2 { // One is a slice, the other is not
		return false, currentPath
	}

	return true, ""
}
