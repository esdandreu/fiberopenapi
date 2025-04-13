package main

import (
	"fmt"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type Parameter struct {
	Name string
	Type string
	// TODO(GIA) Required bool
}

type Operation struct {
	Name        string
	Method      string
	Path        string
	RequestBody string
	Parameters  []Parameter
}

func ExtractOperations(spec *libopenapi.DocumentModel[v3.Document]) []Operation {
	var operations []Operation

	for pair := spec.Model.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
		path := ToFiberPath(pair.Key())
		pathItem := pair.Value()
		parameters := pathItem.Parameters

		if pathItem.Get != nil {
			operations = append(operations,
				extractOperation(path, "Get", parameters, pathItem.Get),
			)
		}
		if pathItem.Put != nil {
			operations = append(operations,
				extractOperation(path, "Put", parameters, pathItem.Put),
			)
		}
		if pathItem.Post != nil {
			operations = append(operations,
				extractOperation(path, "Post", parameters, pathItem.Post),
			)
		}
		if pathItem.Delete != nil {
			operations = append(operations,
				extractOperation(path, "Delete", parameters, pathItem.Delete),
			)
		}
		if pathItem.Options != nil {
			operations = append(operations,
				extractOperation(path, "Options", parameters, pathItem.Options),
			)
		}
		if pathItem.Head != nil {
			operations = append(operations,
				extractOperation(path, "Head", parameters, pathItem.Head),
			)
		}
		if pathItem.Patch != nil {
			operations = append(operations,
				extractOperation(path, "Patch", parameters, pathItem.Patch),
			)
		}
		if pathItem.Trace != nil {
			operations = append(operations,
				extractOperation(path, "Trace", parameters, pathItem.Trace),
			)
		}
	}

	return operations
}

func extractOperation(path, method string, parameters []*v3.Parameter, operation *v3.Operation) Operation {
	if operation.OperationId == "" {
		panic(fmt.Sprintf("operationId is empty for %s %s", method, path))
	}
	result := Operation{
		Name:   ToPascalCase(operation.OperationId),
		Method: method,
		Path:   path,
	}
	if operation.RequestBody != nil {
		result.RequestBody = extractModelFromOperationRequestBody(operation).Name()
	}
	for _, pair := range extractModelsFromOperationParameters(parameters, operation) {
		result.Parameters = append(result.Parameters, Parameter{
			Name: pair.Key(),
			Type: pair.Value().Name(),
		})
	}
	return result
}
