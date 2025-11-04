package browser

import (
	"encoding/json"
	"fmt"
)

// FormatJSON formats any value as indented JSON
func FormatJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}
	return string(b), nil
}

// FormatJSONCompact formats any value as compact JSON
func FormatJSONCompact(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}
	return string(b), nil
}
