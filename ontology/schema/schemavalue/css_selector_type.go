package schemavalue

import (
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type CssSelectorType string

var _ termutil.CustomValue = CssSelectorType("")

func (v CssSelectorType) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.CssSelectorType_Class,
		LexicalForm: string(v),
	}
}

func (CssSelectorType) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v CssSelectorType) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.CssSelectorType_Class {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
