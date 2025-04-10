package main

import (
	"flag"
	"os"
)

func main() {
	var packagePath, outputPath, specPath, typeName string
	flag.StringVar(&packagePath, "path", ".", "path to the package to generate the router for; defaults to current directory")
	flag.StringVar(&outputPath, "output", "handlers.go", "output file name; defaults to handlers.go")
	flag.StringVar(&specPath, "spec", "", "path to the OpenAPI specification file; must be set")
	flag.StringVar(&typeName, "type-name", "Handlers", "name of the interface to generate; defaults to Handlers")
	flag.Parse()
	if specPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	// Load OpenAPI document.
	spec, err := LoadOpenAPIDocument(specPath)
	if err != nil {
		panic(err)
	}

	if err := GenerateHandlers(spec, packagePath, outputPath, typeName); err != nil {
		panic(err)
	}
	if err := GenerateModels(spec, packagePath, "models.go"); err != nil {
		panic(err)
	}
}
