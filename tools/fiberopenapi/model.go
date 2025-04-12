package main

import (
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type ModelType interface {
	Name() string
	Definition() string
}

// Model is a model that has sub-models. Object and array types in the
// OpenAPI specification are modeled as Model. For simplicity, the rest of
// the types are modeled as trivial ModelTrees too.
type Model interface {
	Name() string
	Schema() *base.Schema
	Types() []ModelType
}

// An abstract base model that behaves as a trivial ModelTree.
type baseModel struct {
	name   string
	schema *base.Schema
}

func (m *baseModel) Name() string {
	return m.name
}

func (m *baseModel) Schema() *base.Schema {
	return m.schema
}

type stringModel struct {
	baseModel
}

func (m *stringModel) Definition() string {
	return "string"
}

func (m *stringModel) Types() []ModelType {
	return []ModelType{m}
}

func newStringModel(name string, schema *base.Schema) *stringModel {
	return &stringModel{baseModel{name, schema}}
}

type booleanModel struct {
	baseModel
}

func (m *booleanModel) Definition() string {
	return "bool"
}

func (m *booleanModel) Types() []ModelType {
	return []ModelType{m}
}

func newBooleanModel(name string, schema *base.Schema) *booleanModel {
	return &booleanModel{baseModel{name, schema}}
}

type numberModel struct {
	baseModel
	defaultType string
}

func (m *numberModel) Definition() string {
	switch m.schema.Format {
	case "int64":
		return "int64"
	case "int32":
		return "int32"
	case "float":
		return "float32"
	case "double":
		return "float64"
	default:
		return m.defaultType
	}
}

func (m *numberModel) Types() []ModelType {
	return []ModelType{m}
}

func newNumberModel(name string, schema *base.Schema) *numberModel {
	return &numberModel{baseModel{name, schema}, "float"}
}

func newIntegerModel(name string, schema *base.Schema) *numberModel {
	return &numberModel{baseModel{name, schema}, "int"}
}

type objectModel struct {
	baseModel
	properties []Model
}

func (m *objectModel) Definition() string {
	def := "struct {"
	if len(m.properties) > 0 {
		def += "\n"
	}
	for _, property := range m.properties {
		def += fmt.Sprintf("\t%s %s\n", property.Name(), property.Name())
	}
	def += "}"
	return def
}

func (m *objectModel) Types() []ModelType {
	flattened := []ModelType{m}
	for _, property := range m.properties {
		flattened = append(flattened, property.Types()...)
	}
	return flattened
}

func newObjectModel(name string, schema *base.Schema) *objectModel {
	var properties []Model
	for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
		properties = append(properties, NewModel(pair.Key(), pair.Value()))
	}
	model := &objectModel{baseModel{name, schema}, properties}
	return model
}

// A reference model is a model that doesn't have any model type attached.
type referenceModel struct {
	baseModel
}

func (m *referenceModel) Types() []ModelType {
	return []ModelType{}
}

func newReferenceModel(proxy *base.SchemaProxy) *referenceModel {
	reference := proxy.GetReference()
	if !strings.HasPrefix(reference, "#/components/schemas/") {
		panic(fmt.Errorf("reference not supported: %s", reference))
	}
	name := ToPascalCase(strings.TrimPrefix(reference, "#/components/schemas/"))
	return &referenceModel{baseModel{name, proxy.Schema()}}
}

func NewModel(name string, schemaProxy *base.SchemaProxy) Model {
	if schemaProxy.IsReference() {
		return newReferenceModel(schemaProxy)
	}
	modelName := ToPascalCase(name)
	schema := schemaProxy.Schema()
	schemaType := schema.Type[0]
	switch schemaType {
	case "boolean":
		return newBooleanModel(modelName, schema)
	case "object":
		return newObjectModel(modelName, schema)
	case "number":
		return newNumberModel(modelName, schema)
	case "integer":
		return newIntegerModel(modelName, schema)
	case "string":
		return newStringModel(modelName, schema)
	}
	panic(fmt.Errorf("unsupported type: %s", schemaType))
}

func ExtractModelsFromDocument(spec *libopenapi.DocumentModel[v3.Document]) []ModelType {
	var models []ModelType
	for pair := spec.Model.Components.Schemas.First(); pair != nil; pair = pair.Next() {

	}
	return models
}
