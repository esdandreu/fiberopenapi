package main

import (
	"regexp"
)

var openApiPathParamPattern = regexp.MustCompile(`{([^}]+)}`)

// Convert an OpenAPI path with {param} style parameters to Fiber's :param
// style https://docs.gofiber.io/guide/routing/#parameters.
func ToFiberPath(path string) string {
	return openApiPathParamPattern.ReplaceAllString(path, ":$1")
}
