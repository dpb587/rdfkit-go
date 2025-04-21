package xsdvalue

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type Base64Binary []byte

var _ literalutil.CustomValue = Base64Binary{}

func MapBase64Binary(lexicalForm string) (Base64Binary, error) {
	return Base64Binary(xsdutil.WhiteSpaceCollapse(lexicalForm)), nil
}

func (v Base64Binary) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Base64Binary_Datatype,
		LexicalForm: string(v),
	}
}

func (Base64Binary) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Base64Binary) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Base64Binary_Datatype {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
