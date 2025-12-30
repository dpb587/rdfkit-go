package schematype

import (
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type XPathType string

var _ objecttypes.Value = XPathType("")

func (v XPathType) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.XPathType_Class,
		LexicalForm: string(v),
	}
}

func (XPathType) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v XPathType) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.XPathType_Class {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
