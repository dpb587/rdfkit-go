package jsonldinternal

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
)

type ExpandedIRI interface {
	ExpandedType() string
	Equals(e ExpandedIRI) bool
	NewPropertyValue(propertyOffsets, valueOffsets *cursorio.TextOffsetRange) *ExpandedScalarPrimitive
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

func (e ExpandedIRIasRawValue) NewPropertyValue(propertyOffsets, valueOffsets *cursorio.TextOffsetRange) *ExpandedScalarPrimitive {
	return &ExpandedScalarPrimitive{
		Value:                 e.Value,
		PropertySourceOffsets: propertyOffsets,
	}
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

func (ExpandedIRIasNil) NewPropertyValue(propertyOffsets, valueOffsets *cursorio.TextOffsetRange) *ExpandedScalarPrimitive {
	return &ExpandedScalarPrimitive{
		Value: inspectjson.NullValue{
			SourceOffsets: valueOffsets,
		},
		PropertySourceOffsets: propertyOffsets,
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

func (e ExpandedIRIasIRI) NewPropertyValue(propertyOffsets, valueOffsets *cursorio.TextOffsetRange) *ExpandedScalarPrimitive {
	return &ExpandedScalarPrimitive{
		Value: inspectjson.StringValue{
			SourceOffsets: valueOffsets,
			Value:         string(e),
		},
		PropertySourceOffsets: propertyOffsets,
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

func (e ExpandedIRIasBlankNode) NewPropertyValue(propertyOffsets, valueOffsets *cursorio.TextOffsetRange) *ExpandedScalarPrimitive {
	return &ExpandedScalarPrimitive{
		Value: inspectjson.StringValue{
			SourceOffsets: valueOffsets,
			Value:         string(e),
		},
		PropertySourceOffsets: propertyOffsets,
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

func (e ExpandedIRIasKeyword) NewPropertyValue(propertyOffsets, valueOffsets *cursorio.TextOffsetRange) *ExpandedScalarPrimitive {
	return &ExpandedScalarPrimitive{
		Value: inspectjson.StringValue{
			SourceOffsets: valueOffsets,
			Value:         string(e),
		},
		PropertySourceOffsets: propertyOffsets,
	}
}

func (e ExpandedIRIasKeyword) String() string {
	return string(e)
}
