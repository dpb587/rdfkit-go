package xsdvalue

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type AnyURI string

var _ termutil.CustomValue = AnyURI("")
var _ literalutil.CustomValue = AnyURI("")

func MapAnyURI(lexicalForm string) (AnyURI, error) {
	return AnyURI(xsdutil.WhiteSpaceCollapse(lexicalForm)), nil
}

func (v AnyURI) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v AnyURI) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.AnyURI_Datatype,
		LexicalForm: string(v),
	}
}

func (AnyURI) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v AnyURI) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.AnyURI_Datatype {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
