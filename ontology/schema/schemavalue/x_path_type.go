package schemavalue

import (
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type XPathType string

var _ termutil.CustomValue = XPathType("")

func (v XPathType) AsTerm() rdf.Term {
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
