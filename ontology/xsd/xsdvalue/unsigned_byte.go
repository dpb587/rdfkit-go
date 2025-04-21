package xsdvalue

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type UnsignedByte uint8

var _ literalutil.CustomValue = UnsignedByte(0)

func MapUnsignedByte(lexicalForm string) (UnsignedByte, error) {
	vInt64, err := strconv.ParseUint(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 8)
	if err != nil {
		return UnsignedByte(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return UnsignedByte(vInt64), nil
}

func (v UnsignedByte) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.UnsignedByte_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (UnsignedByte) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v UnsignedByte) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.UnsignedByte_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
