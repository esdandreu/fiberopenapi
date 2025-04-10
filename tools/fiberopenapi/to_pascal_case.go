package main

import (
	"regexp"
	"strings"
)

// Converts a string from kebab-case, snake_case, or camelCase to PascalCase
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// First replace hyphens with underscores to handle both kebab-case and snake_case
	s = strings.ReplaceAll(s, "-", "_")

	// Handle camelCase by inserting underscores before capital letters
	// This regex finds positions before capital letters that are not at the start of the string
	// and not already preceded by an underscore
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	s = re.ReplaceAllString(s, "${1}_${2}")

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
