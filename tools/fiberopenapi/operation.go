package main

import (
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type Parameter struct {
	Name string
	Type string
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
		result.RequestBody = parseRequestBodyName(result.Name, operation.RequestBody)
	}
	result.Parameters = make([]Parameter, len(parameters)+len(operation.Parameters))
	for i, parameter := range parameters {
		result.Parameters[i] = parseParameter(result.Name, parameter)
	}
	for i, parameter := range operation.Parameters {
		result.Parameters[len(parameters)+i] = parseParameter(result.Name, parameter)
	}
	return result
}

func getReferenceName(proxy *base.SchemaProxy) string {
	return ToPascalCase(strings.TrimPrefix(proxy.GetReference(), "#/components/schemas/"))
}

func parseRequestBodyName(parentName string, requestBody *v3.RequestBody) string {
	content := requestBody.Content.GetOrZero("application/json")
	if content == nil {
		panic(fmt.Sprintf("no JSON content for %s request body", parentName))
	}
	proxy := content.Schema
	if proxy.IsReference() {
		return getReferenceName(proxy)
	}
	return parentName + "Item"
}

func parseParameter(parentName string, parameter *v3.Parameter) Parameter {
	proxy := parameter.Schema
	if proxy.IsReference() {
		return Parameter{
			Name: parameter.Name,
			Type: getReferenceName(proxy),
		}
	}
	return Parameter{
		Name: parameter.Name,
		Type: parentName + ToPascalCase(parameter.Name),
	}
}
