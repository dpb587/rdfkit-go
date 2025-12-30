package schematype

import (
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type Integer int64

var _ objecttypes.Value = Integer(0)

func (v Integer) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.Integer_Class,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (Integer) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Integer) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.Integer_Class {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
