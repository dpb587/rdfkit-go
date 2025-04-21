package jsonldinternal

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
)

type ExpandedIRI interface {
	ExpandedType() string
	Equals(e ExpandedIRI) bool
	NewValue(offsetRange *cursorio.TextOffsetRange) inspectjson.Value
	String() string
}

// [dpb] weird; change abstraction?

type ExpandedIRIasRawValue struct {
	Value inspectjson.Value
}

func (ExpandedIRIasRawValue) ExpandedType() string {
	return "value"
}

func (e ExpandedIRIasRawValue) Equals(e2 ExpandedIRI) bool {
	ee2, ok := e2.(ExpandedIRIasRawValue)
	return ok && ee2.Value == e.Value
}

func (e ExpandedIRIasRawValue) NewValue(offsetRange *cursorio.TextOffsetRange) inspectjson.Value {
	return e.Value
}

func (ExpandedIRIasRawValue) String() string {
	return "null" // not actually valid
}

//

type ExpandedIRIasNil struct{}

func (ExpandedIRIasNil) ExpandedType() string {
	return "null"
}

func (ExpandedIRIasNil) Equals(e ExpandedIRI) bool {
	_, ok := e.(ExpandedIRIasNil)
	return ok
}

func (ExpandedIRIasNil) NewValue(offsetRange *cursorio.TextOffsetRange) inspectjson.Value {
	return inspectjson.NullValue{
		SourceOffsets: offsetRange,
	}
}

func (ExpandedIRIasNil) String() string {
	return "null" // not actually valid
}

type ExpandedIRIasIRI string

func (ExpandedIRIasIRI) ExpandedType() string {
	return "iri"
}

func (e ExpandedIRIasIRI) Equals(e2 ExpandedIRI) bool {
	e2IRI, ok := e2.(ExpandedIRIasIRI)

	return ok && e == e2IRI
}

func (e ExpandedIRIasIRI) NewValue(offsetRange *cursorio.TextOffsetRange) inspectjson.Value {
	return inspectjson.StringValue{
		SourceOffsets: offsetRange,
		Value:         string(e),
	}
}

func (e ExpandedIRIasIRI) String() string {
	return string(e)
}

type ExpandedIRIasBlankNode string

func (ExpandedIRIasBlankNode) ExpandedType() string {
	return "blank node"
}

func (e ExpandedIRIasBlankNode) Equals(e2 ExpandedIRI) bool {
	e2t, ok := e2.(ExpandedIRIasBlankNode)

	return ok && e == e2t
}

func (e ExpandedIRIasBlankNode) NewValue(offsetRange *cursorio.TextOffsetRange) inspectjson.Value {
	return inspectjson.StringValue{
		SourceOffsets: offsetRange,
		Value:         string(e),
	}
}

func (e ExpandedIRIasBlankNode) String() string {
	return string(e)
}

type ExpandedIRIasKeyword string

func (ExpandedIRIasKeyword) ExpandedType() string {
	return "keyword"
}

func (e ExpandedIRIasKeyword) Equals(e2 ExpandedIRI) bool {
	e2t, ok := e2.(ExpandedIRIasKeyword)

	return ok && string(e) == string(e2t)
}

func (e ExpandedIRIasKeyword) NewValue(offsetRange *cursorio.TextOffsetRange) inspectjson.Value {
	return inspectjson.StringValue{
		SourceOffsets: offsetRange,
		Value:         string(e),
	}
}

func (e ExpandedIRIasKeyword) String() string {
	return string(e)
}
