package xsdvalue

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type HexBinary []byte

var _ literalutil.CustomValue = HexBinary{}

func MapHexBinary(lexicalForm string) (HexBinary, error) {
	return HexBinary(xsdutil.WhiteSpaceCollapse(lexicalForm)), nil
}

func (v HexBinary) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.HexBinary_Datatype,
		LexicalForm: string(v),
	}
}

func (HexBinary) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v HexBinary) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.HexBinary_Datatype {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
