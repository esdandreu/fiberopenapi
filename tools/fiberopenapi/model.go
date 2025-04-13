package main

import (
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type ModelType interface {
	Name() string
	Docstring() string
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

func (m *baseModel) Docstring() string {
	if m.schema.Deprecated != nil && *m.schema.Deprecated {
		if m.schema.Description != "" {
			return fmt.Sprintf("// Deprecated: %s\n", m.schema.Description)
		} else {
			return "// Deprecated\n"
		}
	}
	if m.schema.Description != "" {
		return fmt.Sprintf("// %s\n", m.schema.Description)
	}
	// ? schema.ExternalDocs
	// ? schema.Example
	// ? schema.Examples
	return ""
}

type nullModel struct {
	baseModel
}

func (m *nullModel) Definition() string {
	return "Null"
}

func (m *nullModel) Types() []ModelType {
	return []ModelType{m}
}

type stringModel struct {
	baseModel
}

func (m *stringModel) Definition() string {
	// TODO(GIA) Time formats: date-time, date, time etc.
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
	properties *orderedmap.Map[string, Model]
}

func (m *objectModel) Definition() string {
	var def string
	for pair := m.properties.First(); pair != nil; pair = pair.Next() {
		property := pair.Value()
		// TODO(GIA) Apply nullable and required
		def += fmt.Sprintf("\t%s %s `json:\"%s,omitempty\"`\n",
			property.Name(), property.Name(), pair.Key(),
		)
	}
	if def != "" {
		return "struct {\n" + def + "}"
	}
	return "struct {}"
}

func (m *objectModel) Types() []ModelType {
	flattened := []ModelType{m}
	for property := range m.properties.ValuesFromOldest() {
		flattened = append(flattened, property.Types()...)
	}
	return flattened
}

func newObjectModel(name string, schema *base.Schema) *objectModel {
	properties := orderedmap.New[string, Model]()
	for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
		properties.Set(pair.Key(), NewModel(pair.Key(), pair.Value()))
	}
	model := &objectModel{baseModel{name, schema}, properties}
	return model
}

type arrayModel struct {
	baseModel
	items Model
}

func (m *arrayModel) Definition() string {
	return fmt.Sprintf("[]%s", m.items.Name())
}

func (m *arrayModel) Types() []ModelType {
	return append([]ModelType{m}, m.items.Types()...)
}

func newArrayModel(name string, schema *base.Schema) *arrayModel {
	if schema.Items == nil {
		panic(fmt.Errorf("array type must have an items property"))
	}
	if schema.Items.IsB() {
		panic(fmt.Errorf("array type with boolean items is not supported"))
	}
	items := NewModel(name+"Item", schema.Items.A)
	return &arrayModel{baseModel{name, schema}, items}
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

// TODO(GIA) unionModel

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
	case "array":
		return newArrayModel(modelName, schema)
	case "number":
		return newNumberModel(modelName, schema)
	case "integer":
		return newIntegerModel(modelName, schema)
	case "string":
		return newStringModel(modelName, schema)
	}
	panic(fmt.Errorf("unsupported type: %s", schemaType))
}

func ExtractModelTypesFromDocument(spec *libopenapi.DocumentModel[v3.Document]) []ModelType {
	if spec == nil {
		return nil
	}
	// Extract the model tree from the document (each model can have
	// nested models).
	var models []Model
	models = append(models, extractModelsFromComponents(spec.Model.Components)...)
	models = append(models, extractModelsFromPaths(spec.Model.Paths)...)
	// Flatten the models into a single slice of model types.
	var modelTypes []ModelType
	for _, model := range models {
		modelTypes = append(modelTypes, model.Types()...)
	}
	return modelTypes
}

func extractModelsFromComponents(components *v3.Components) []Model {
	if components == nil {
		return nil
	}
	var models []Model
	for pair := components.Schemas.First(); pair != nil; pair = pair.Next() {
		models = append(models, NewModel(pair.Key(), pair.Value()))
	}
	return models
}

func extractModelsFromPaths(paths *v3.Paths) []Model {
	if paths == nil {
		return nil
	}
	var models []Model
	for pair := paths.PathItems.First(); pair != nil; pair = pair.Next() {
		pathItem := pair.Value()
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Get)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Put)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Post)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Delete)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Options)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Head)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Patch)...)
		models = append(models, extractModelsFromOperation(pathItem.Parameters, pathItem.Trace)...)
	}
	return models
}

func extractModelsFromOperation(
	pathItemParameters []*v3.Parameter, operation *v3.Operation,
) []Model {
	if operation == nil {
		return nil
	}
	var models []Model
	if operation.RequestBody != nil {
		models = append(models, extractModelFromOperationRequestBody(operation))
	}
	for _, pair := range extractModelsFromOperationParameters(
		pathItemParameters, operation,
	) {
		models = append(models, pair.Value())
	}
	return models
}

func extractModelFromOperationRequestBody(operation *v3.Operation) Model {
	content := operation.RequestBody.Content.GetOrZero("application/json")
	if content == nil {
		panic(fmt.Sprintf("no JSON content for %s request body", operation.OperationId))
	}
	return NewModel(operation.OperationId+"RequestBody", content.Schema)
}

func extractModelsFromOperationParameters(
	pathItemParameters []*v3.Parameter, operation *v3.Operation,
) []orderedmap.Pair[string, Model] {
	return append(
		extractModelsFromParameters(operation.OperationId, operation.Parameters),
		extractModelsFromParameters(operation.OperationId, pathItemParameters)...,
	)
}

func extractModelsFromParameters(prefix string, parameters []*v3.Parameter) []orderedmap.Pair[string, Model] {
	models := make([]orderedmap.Pair[string, Model], len(parameters))
	for i, parameter := range parameters {
		models[i] = orderedmap.NewPair(
			ToCamelCase(parameter.Name),
			NewModel(prefix+ToPascalCase(parameter.Name), parameter.Schema),
		)
	}
	return models
}
