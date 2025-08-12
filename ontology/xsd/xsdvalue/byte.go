package xsdvalue

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Byte int8

var _ termutil.CustomValue = Byte(0)
var _ literalutil.CustomValue = Byte(0)

func MapByte(lexicalForm string) (Byte, error) {
	vInt64, err := strconv.ParseInt(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 8)
	if err != nil {
		return Byte(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Byte(vInt64), nil
}

func (v Byte) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v Byte) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Byte_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (Byte) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Byte) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Byte_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
