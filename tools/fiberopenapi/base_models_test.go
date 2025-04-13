package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Task struct {
	Error *Error `json:"error,omitempty"`
}

// Model generates structs like this for union types.
type Error struct {
	value any
}

func (e *Error) Null() (Null, bool) {
	s, ok := e.value.(Null)
	return s, ok
}

func (e *Error) IsNull() bool {
	_, ok := e.Null()
	return ok
}

func NewErrorAsNull() *Error {
	return &Error{value: Null{}}
}

func (e *Error) String() (string, bool) {
	s, ok := e.value.(string)
	return s, ok
}

func (e *Error) IsString() bool {
	_, ok := e.String()
	return ok
}

func NewErrorAsString(s string) *Error {
	return &Error{value: s}
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var errs []error
	var null Null
	if err := json.Unmarshal(data, &null); err == nil {
		e.value = null
		return nil
	} else {
		errs = append(errs, err)
	}
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		e.value = str
		return nil
	} else {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func TestTaskUnmarshal(t *testing.T) {
	testCases := map[string]struct {
		json        string
		expected    Task
		expectedErr bool
	}{
		"with null": {
			json:     `{"error": null}`,
			expected: Task{Error: NewErrorAsNull()},
		},
		"with string": {
			json:     `{"error": "something went wrong"}`,
			expected: Task{Error: NewErrorAsString("something went wrong")},
		},
		"with empty string": {
			json:     `{"error": ""}`,
			expected: Task{Error: NewErrorAsString("")},
		},
		"empty object": {
			json:     `{}`,
			expected: Task{Error: nil},
		},
		"invalid error type": {
			json:        `{"error": 123}`,
			expectedErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var task Task
			err := json.Unmarshal([]byte(tc.json), &task)
			if tc.expectedErr {
				t.Logf("error: %v", err)
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, task)
			}
		})
	}
}
