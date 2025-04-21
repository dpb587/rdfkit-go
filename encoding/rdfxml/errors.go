package rdfxml

import (
	"encoding/xml"
	"errors"
	"fmt"
)

var (
	ErrDirectivesNotSupported = errors.New("directives not supported")
	ErrAttributeNotAllowed    = errors.New("attribute not allowed")
	ErrElementNotAllowed      = errors.New("element not allowed")
)

//

type AttributeNotAllowedError struct {
	Name xml.Name
}

func (e AttributeNotAllowedError) Error() string {
	return fmt.Sprintf("%s: %s (%s)", ErrAttributeNotAllowed, e.Name.Local, e.Name.Space)
}

//

type ElementNotAllowedError struct {
	Name xml.Name
}

func (e ElementNotAllowedError) Error() string {
	return fmt.Sprintf("%s: %s (%s)", ErrElementNotAllowed, e.Name.Local, e.Name.Space)
}

//

type DuplicateScopedNameError struct {
	Name string
}

func (e DuplicateScopedNameError) Error() string {
	return fmt.Sprintf("duplicate name found in scope: %s", e.Name)
}

//

type InvalidNameError struct {
	Name string
}

func (e InvalidNameError) Error() string {
	return fmt.Sprintf("invalid name: %s", e.Name)
}
