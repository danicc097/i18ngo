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

	for i := 1; i < len(structures); i++ {
		if ok := compareMaps(structures[0], structures[i]); !ok {
			return fmt.Errorf("structure mismatch between translation files")
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

func compareMaps(map1, map2 anyMap) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key1, val1 := range map1 {
		val2, exists := map2[key1]
		if !exists {
			return false
		}

		if ok := compareValues(val1, val2); !ok {
			return false
		}
	}

	return true
}

func compareValues(val1, val2 interface{}) bool {
	map1, ok1 := val1.(anyMap)
	map2, ok2 := val2.(anyMap)
	if ok1 && ok2 {
		return compareMaps(map1, map2)
	}

	_, ok1 = val1.([]interface{})
	_, ok2 = val2.([]interface{})
	if ok1 && ok2 || !ok1 && !ok2 {
		return true // dont care about values, just that both are either slices or neither are
	}

	// types don't match, assume its an error
	return false
}
