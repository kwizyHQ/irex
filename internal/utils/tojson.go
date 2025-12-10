package utils

import (
	"encoding/json"
)

// ToJSON converts any Go value to a pretty-printed JSON string
func ToJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
