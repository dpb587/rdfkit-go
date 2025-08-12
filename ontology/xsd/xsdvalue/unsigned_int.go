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

type UnsignedInt uint32

var _ termutil.CustomValue = UnsignedInt(0)
var _ literalutil.CustomValue = UnsignedInt(0)

func MapUnsignedInt(lexicalForm string) (UnsignedInt, error) {
	vInt64, err := strconv.ParseUint(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 32)
	if err != nil {
		return UnsignedInt(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return UnsignedInt(vInt64), nil
}

func (v UnsignedInt) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v UnsignedInt) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.UnsignedInt_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (UnsignedInt) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v UnsignedInt) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.UnsignedInt_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
