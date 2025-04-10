package main

import (
	"testing"
)

func TestToPascalCase(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
	}{
		"kebab-case to PascalCase": {
			input:    "hello-world",
			expected: "HelloWorld",
		},
		"snake_case to PascalCase": {
			input:    "hello_world",
			expected: "HelloWorld",
		},
		"camelCase to PascalCase": {
			input:    "helloWorld",
			expected: "HelloWorld",
		},
		"mixed case with hyphens and underscores": {
			input:    "hello-world_example",
			expected: "HelloWorldExample",
		},
		"empty string": {
			input:    "",
			expected: "",
		},
		"single word": {
			input:    "hello",
			expected: "Hello",
		},
		"multiple consecutive separators": {
			input:    "hello--world__example",
			expected: "HelloWorldExample",
		},
		"already in PascalCase": {
			input:    "HelloWorld",
			expected: "HelloWorld",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := ToPascalCase(tc.input)
			if result != tc.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}
