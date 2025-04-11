package main

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Read and parse an OpenAPI specification file.
func LoadOpenAPIDocument(specPath string) (*libopenapi.DocumentModel[v3.Document], error) {
	specByteArray, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read specification file at %s: %w", specPath, err)
	}
	return loadOpenAPIDocument(specByteArray)
}

func loadOpenAPIDocument(specByteArray []byte) (*libopenapi.DocumentModel[v3.Document], error) {
	document, err := libopenapi.NewDocument(specByteArray)
	if err != nil {
		return nil, fmt.Errorf("cannot create new document: %w", err)
	}
	model, errors := document.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %v\n", errors[i])
		}
		return nil, fmt.Errorf("cannot create v3 model from document: %d errors reported", len(errors))
	}
	return model, nil
}
