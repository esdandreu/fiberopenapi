package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUndefined(t *testing.T) {
	testCases := map[string]struct {
		input    any
		expected bool
	}{
		"nil":     {nil, true},
		"not nil": {1, false},
		"string":  {"hello", false},
		"nil string pointer": {
			func() *string {
				var s *string
				return s
			}(),
			true,
		},
		"non nil string pointer": {
			func() *string {
				s := "hello"
				return &s
			}(),
			false,
		},
		"empty array": {
			[]string{},
			false,
		},
		"nil array": {
			func() []string {
				return nil
			}(),
			true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsUndefined(tc.input))
		})
	}
}

func TestIsNull(t *testing.T) {
	testCases := map[string]struct {
		input    any
		expected bool
	}{
		"nil":          {nil, false},
		"not nullable": {1, false},
		"string":       {"hello", false},
		"nullable with null true": {
			Nullable[string]{
				Value:  "hello",
				isNull: true,
			},
			true,
		},
		"nullable with null false": {
			Nullable[string]{
				Value:  "hello",
				isNull: false,
			},
			false,
		},
		"nullable with nil value": {
			Nullable[*string]{
				Value:  nil,
				isNull: false,
			},
			false,
		},
		"nullable with nil value and isNull true": {
			Nullable[*string]{
				Value:  nil,
				isNull: true,
			},
			true,
		},
		// "NewNull string": {
		// 	Null[Nullable[string]](),
		// 	true,
		// },
		// "NewNull *string": {
		// 	Null[Nullable[*string]](),
		// 	true,
		// },
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsNull(tc.input))
		})
	}
}
