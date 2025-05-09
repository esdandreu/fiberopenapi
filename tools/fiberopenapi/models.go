package main

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrMaxLength = errors.New("maxLength")
	ErrMinLength = errors.New("minLength")
)

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

func NewValidationError(format string, a ...any) error {
	return ValidationError{Err: fmt.Errorf(format, a...)}
}

func NewMaxLengthError(got int, want int) error {
	return NewValidationError("%w: got %d, want %d", ErrMaxLength, got, want)
}

func NewMinLengthError(got int, want int) error {
	return NewValidationError("%w: got %d, want %d", ErrMinLength, got, want)
}

type Nullable[T any] struct {
	Value  T
	isNull bool
}

func (n Nullable[T]) IsNull() bool {
	return n.isNull
}

type nullable interface {
	IsNull() bool
}

func IsNull(v any) bool {
	if v == nil {
		return false
	}
	if n, ok := v.(nullable); ok {
		return n.IsNull()
	}
	return false
}

// func Null[T any]() T {
// 	return T{isNull: true}
// }

func IsUndefined(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan,
		reflect.Func, reflect.Interface:
		return rv.IsNil()
	}
	return false
}
