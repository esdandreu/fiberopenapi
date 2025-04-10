package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
)

// Contains the generated code. One can append content to it using Printf. Once
// done, one can use WriteFile to write formatted content to a file.
type Generator struct {
	content bytes.Buffer
}

func (g *Generator) Printf(format string, args ...any) {
	fmt.Fprintf(&g.content, format, args...)
}

func (g *Generator) WriteFile(path string) error {
	src, err := format.Source(g.content.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}
	err = os.WriteFile(path, src, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write generated code: %w", err)
	}
	return nil
}
