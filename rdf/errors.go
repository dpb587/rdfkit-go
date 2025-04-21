package rdf

import "errors"

var (
	ErrLiteralDatatypeNotValid    = errors.New("literal datatype not valid")
	ErrLiteralLexicalFormNotValid = errors.New("literal lexical form not valid")
)
