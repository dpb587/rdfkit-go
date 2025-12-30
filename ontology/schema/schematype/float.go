package schematype

import (
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type Float float64

var _ objecttypes.Value = Float(0)

func (v Float) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.Float_Class,
		LexicalForm: strconv.FormatFloat(float64(v), 'f', -1, 64),
	}
}

func (Float) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Float) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.Float_Class {
		return false
	}

	return strconv.FormatFloat(float64(v), 'f', -1, 64) == tLiteral.LexicalForm
}
