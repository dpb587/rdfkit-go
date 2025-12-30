package rdftype

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type HTML string

var _ objecttypes.Value = HTML("")

func (v HTML) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    rdfiri.HTML_Datatype,
		LexicalForm: string(v),
	}
}

func (HTML) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v HTML) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != rdfiri.HTML_Datatype {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
