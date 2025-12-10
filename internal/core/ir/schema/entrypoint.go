package schema

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/utils"
)

// ParseFile parses a single HCL file and returns a Model slice (single element)
func ParseFile(path string) ([]Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	models, err := parseHCLFile(path)
	if err != nil {
		return nil, fmt.Errorf("parse file %s: %w", path, err)
	}
	return models, nil
}

// ParseDir finds all .hcl files under dir and merges parsed models into a single slice
func ParseDir(dir string) ([]Model, error) {
	var out []Model
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".hcl" {
			return nil
		}
		ms, err := parseHCLFile(p)
		if err != nil {
			return fmt.Errorf("parse %s: %w", p, err)
		}
		out = append(out, ms...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetJsonModels accepts either a file path or directory path. If path is a directory
// it parses all .hcl files under it and returns a JSON representation of the merged models.
func GetJsonModels(path string) (string, error) {
	// determine if path is file or dir
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat path: %w", err)
	}

	var models []Model
	if info.IsDir() {
		models, err = ParseDir(path)
		if err != nil {
			return "", err
		}
	} else {
		models, err = ParseFile(path)
		if err != nil {
			return "", err
		}
	}

	s, err := utils.ToJSON(models)
	if err != nil {
		return "", fmt.Errorf("json encode: %w", err)
	}
	return s, nil
}
