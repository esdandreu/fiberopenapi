package main

import (
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Route represents an API endpoint with its method, path and operation ID
type Route struct {
	Method      string
	Path        string
	OperationId string
}

// ExtractRoutes returns a list of routes from an OpenAPI specification
func ExtractRoutes(spec *libopenapi.DocumentModel[v3.Document]) []Route {
	var routes []Route

	for pair := spec.Model.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
		path := ToFiberPath(pair.Key())
		pathItem := pair.Value()

		if pathItem.Get != nil && pathItem.Get.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Get",
				Path:        path,
				OperationId: pathItem.Get.OperationId,
			})
		}
		if pathItem.Put != nil && pathItem.Put.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Put",
				Path:        path,
				OperationId: pathItem.Put.OperationId,
			})
		}
		if pathItem.Post != nil && pathItem.Post.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Post",
				Path:        path,
				OperationId: pathItem.Post.OperationId,
			})
		}
		if pathItem.Delete != nil && pathItem.Delete.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Delete",
				Path:        path,
				OperationId: pathItem.Delete.OperationId,
			})
		}
		if pathItem.Options != nil && pathItem.Options.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Options",
				Path:        path,
				OperationId: pathItem.Options.OperationId,
			})
		}
		if pathItem.Head != nil && pathItem.Head.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Head",
				Path:        path,
				OperationId: pathItem.Head.OperationId,
			})
		}
		if pathItem.Patch != nil && pathItem.Patch.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Patch",
				Path:        path,
				OperationId: pathItem.Patch.OperationId,
			})
		}
		if pathItem.Trace != nil && pathItem.Trace.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Trace",
				Path:        path,
				OperationId: pathItem.Trace.OperationId,
			})
		}
	}

	return routes
}
