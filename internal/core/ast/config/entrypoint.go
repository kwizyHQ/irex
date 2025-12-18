package config

import (
	"fmt"

	"github.com/kwizyHQ/irex/internal/utils"
)

// Parse parses the HCL file at the given path and returns the ConfigDefinition struct
func Parse(path string) (*ConfigDefinition, error) {
	def, err := parseHCLFile(path)
	if err != nil {
		return nil, err
	}
	return def, nil
}

// GetJson parses the HCL file and returns its JSON representation
func GetJson(path string) (string, error) {
	def, err := Parse(path)
	if err != nil {
		return "", err
	}
	b, err := utils.ToJSON(def)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(b), nil
}
