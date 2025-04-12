package main

import (
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type Schema struct {
	*base.Schema
	// The name of the schema in the OpenAPI document.
	Name string
}

// TODO(GIA) Recursive when SchemaProxy is not a reference.
func NewSchemas(name string, proxy *base.SchemaProxy) []*Schema {
	return nil
}

func ExtractSchemasFromDocument(spec *libopenapi.DocumentModel[v3.Document]) []*Schema {
	var schemas []*Schema
	for pair := spec.Model.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		// TODO(GIA) Recursive
		schemas = append(schemas, &Schema{pair.Value().Schema(), pair.Key()})
	}
	return schemas
}

func ExtractSchemasFromProperties(properties *orderedmap.Map[string, *base.SchemaProxy]) []*Schema {
	var schemas []*Schema
	for pair := properties.First(); pair != nil; pair = pair.Next() {
		schemas = append(schemas, &Schema{pair.Value().Schema(), pair.Key()})
	}
	return schemas
}

// TODO(GIA) Properties?

// A pascal-cased version of the schema name.
func (s *Schema) ModelName() string {
	return ToPascalCase(s.Name)
}

// The name of the model type. Usually the same as the model name.
func (s *Schema) ModelType() string {
	modelName := s.ModelName()
	if s.Nullable != nil && *s.Nullable {
		return "Nullable[" + modelName + "]"
	}
	return modelName
}

// TODO(GIA) Go Type?
