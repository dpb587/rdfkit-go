package xsdtype

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type String string

var _ objecttypes.Value = String("")

func MapString(lexicalForm string) (String, error) {
	return String(lexicalForm), nil
}

func (v String) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    xsdiri.String_Datatype,
		LexicalForm: string(v),
	}
}

func (String) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v String) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.String_Datatype {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
