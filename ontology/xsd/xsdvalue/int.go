package xsdvalue

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type Int int32

var _ literalutil.CustomValue = Int(0)

func MapInt(lexicalForm string) (Int, error) {
	vInt64, err := strconv.ParseInt(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 32)
	if err != nil {
		return Int(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Int(vInt64), nil
}

func (v Int) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Int_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (Int) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Int) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Int_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
