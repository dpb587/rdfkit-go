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

type UnsignedShort uint16

var _ termutil.CustomValue = UnsignedShort(0)
var _ literalutil.CustomValue = UnsignedShort(0)

func MapUnsignedShort(lexicalForm string) (UnsignedShort, error) {
	vInt64, err := strconv.ParseUint(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 16)
	if err != nil {
		return UnsignedShort(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return UnsignedShort(vInt64), nil
}

func (v UnsignedShort) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v UnsignedShort) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.UnsignedShort_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (UnsignedShort) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v UnsignedShort) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.UnsignedShort_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
