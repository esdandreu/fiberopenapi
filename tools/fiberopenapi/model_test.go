package main

import (
	"testing"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/assert"
)

func TestObjectModel(t *testing.T) {
	testCases := map[string]struct {
		schema                  *base.Schema
		expectedTypeDefinitions map[string]string
	}{
		"Status": {
			schema: func() *base.Schema {
				schema := &base.Schema{Type: []string{"object"}}
				schema.Properties = orderedmap.New[string, *base.SchemaProxy]()
				schema.Properties.Set("winner", base.CreateSchemaProxy(
					&base.Schema{Type: []string{"string"}},
				))
				schema.Properties.Set("board", base.CreateSchemaProxy(
					&base.Schema{Type: []string{"string"}},
				))
				return schema
			}(),
			expectedTypeDefinitions: map[string]string{
				"Status": `struct {
	Winner Winner
	Board Board
}`,
				"Winner": "string",
				"Board":  "string",
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			model := newObjectModel(name, testCase.schema)
			assert.Equal(t, name, model.Name())
			typeDefinitions := map[string]string{}
			for _, modelType := range model.Types() {
				typeDefinitions[modelType.Name()] = modelType.Definition()
			}
			assert.Equal(t, testCase.expectedTypeDefinitions, typeDefinitions)
		})
	}
}
