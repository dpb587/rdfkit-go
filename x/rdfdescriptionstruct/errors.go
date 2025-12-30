package rdfdescriptionstruct

import "fmt"

// UnmarshalError represents an error during unmarshaling.
type UnmarshalError struct {
	Field string
	Err   error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal error at field %q: %v", e.Field, e.Err)
}

func (e *UnmarshalError) Unwrap() error {
	return e.Err
}

// TypeMismatchError represents a type mismatch during unmarshaling.
type TypeMismatchError struct {
	Expected string
	Got      string
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch: expected %s, got %s", e.Expected, e.Got)
}

// InvalidTagError represents an invalid struct tag.
type InvalidTagError struct {
	Tag string
	Err error
}

func (e *InvalidTagError) Error() string {
	return fmt.Sprintf("invalid tag %q: %v", e.Tag, e.Err)
}

func (e *InvalidTagError) Unwrap() error {
	return e.Err
}
