package main

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

// Loads the package name of a directory.
func LoadPackageName(dir string) (string, error) {
	cfg := &packages.Config{Mode: packages.NeedName}
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return "", fmt.Errorf("cannot load package info: %w", err)
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found")
	}
	return pkgs[0].Name, nil
}
