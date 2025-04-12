package main

import (
	"regexp"
	"strings"
)

// Converts a string from kebab-case, snake_case, or camelCase to PascalCase.
// It also handles OpenAPI /path/{param} paths.
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Replace all non-alphanumeric characters with underscores
	nonAlphanumericPattern := regexp.MustCompile(`[^a-zA-Z0-9]`)
	s = nonAlphanumericPattern.ReplaceAllString(s, "_")

	// Handle camelCase by inserting underscores before capital letters
	// This regex finds positions before capital letters that are not at the start of the string
	// and not already preceded by an underscore
	camelCasePattern := regexp.MustCompile(`([a-z])([A-Z])`)
	s = camelCasePattern.ReplaceAllString(s, "${1}_${2}")

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
