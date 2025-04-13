package main

import "strings"

func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}
	pascalCase := ToPascalCase(s)
	return strings.ToLower(pascalCase[:1]) + pascalCase[1:]
}
