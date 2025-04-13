package main

import (
	"testing"
)

func TestToCamelCase(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
	}{
		"kebab-case to camelCase": {
			input:    "hello-world",
			expected: "helloWorld",
		},
		"snake_case to camelCase": {
			input:    "hello_world",
			expected: "helloWorld",
		},
		"PascalCase to camelCase": {
			input:    "HelloWorld",
			expected: "helloWorld",
		},
		"mixed case with hyphens and underscores": {
			input:    "hello-world_example",
			expected: "helloWorldExample",
		},
		"empty string": {
			input:    "",
			expected: "",
		},
		"single word": {
			input:    "hello",
			expected: "hello",
		},
		"multiple consecutive separators": {
			input:    "hello--world__example",
			expected: "helloWorldExample",
		},
		"already in camelCase": {
			input:    "helloWorld",
			expected: "helloWorld",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := ToCamelCase(tc.input)
			if result != tc.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}
