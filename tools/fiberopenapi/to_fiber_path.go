package main

import (
	"regexp"
)

// Convert an OpenAPI path with {param} style parameters to Fiber's :param
// style https://docs.gofiber.io/guide/routing/#parameters.
func ToFiberPath(path string) string {
	re := regexp.MustCompile(`{([^}]+)}`)
	return re.ReplaceAllString(path, ":$1")
}
