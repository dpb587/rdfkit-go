package iri

import "fmt"

var (
	errUnknownPrefix = "unknown prefix"
)

type UnknownPrefixError struct {
	Prefix string
}

func NewUnknownPrefixError(prefix string) UnknownPrefixError {
	return UnknownPrefixError{
		Prefix: prefix,
	}
}

func (e UnknownPrefixError) Error() string {
	return fmt.Sprintf("%s: %s", errUnknownPrefix, e.Prefix)
}
