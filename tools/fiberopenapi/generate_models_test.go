package main

import (
	"testing"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateModel(t *testing.T) {
	testCases := map[string]struct {
		schema   *base.Schema
		expected string
	}{
		"errorMessage": {
			schema: func() *base.Schema {
				schema := &base.Schema{Type: []string{"string"}}
				maxLength := int64(256)
				schema.MaxLength = &maxLength
				return schema
			}(),
			expected: `
type ErrorMessage string

func (v ErrorMessage) Validate() error {
	var errs []error
	if got := len(v); got > 256 {
		errs = append(errs, NewMaxLengthError(got, 256))
	}
	return errors.Join(errs...)
}
`,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			g := &Generator{}
			require.NoError(t, generateModel(g, &Schema{testCase.schema, name}))
			assert.Equal(t, testCase.expected, g.content.String())
		})
	}
}
