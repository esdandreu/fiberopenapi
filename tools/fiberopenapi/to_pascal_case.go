package main

import (
	"strings"
)

// TODO(GIA) Handle camelCase.
// Converts a string from kebab-case or snake_case to PascalCase
func ToPascalCase(s string) string {
	// First replace hyphens with underscores to handle both kebab-case and snake_case
	s = strings.ReplaceAll(s, "-", "_")

	// Split the string by underscores
	parts := strings.Split(s, "_")

	// Convert each part to title case and join them
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, "")
}
