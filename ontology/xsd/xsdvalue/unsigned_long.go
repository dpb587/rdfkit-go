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

type UnsignedLong uint64

var _ termutil.CustomValue = UnsignedLong(0)
var _ literalutil.CustomValue = UnsignedLong(0)

func MapUnsignedLong(lexicalForm string) (UnsignedLong, error) {
	vInt64, err := strconv.ParseUint(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 64)
	if err != nil {
		return UnsignedLong(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return UnsignedLong(vInt64), nil
}

func (v UnsignedLong) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v UnsignedLong) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.UnsignedLong_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (UnsignedLong) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v UnsignedLong) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.UnsignedLong_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
