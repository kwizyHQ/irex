package service

import (
	"encoding/json"
	"fmt"
)

// Parse parses the HCL file at the given path and returns the ServiceDefinition struct
func Parse(path string) (*ServiceDefinition, error) {
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
	b, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(b), nil
}
