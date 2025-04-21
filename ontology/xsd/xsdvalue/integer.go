package xsdvalue

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type Integer int64

var _ literalutil.CustomValue = Integer(0)

func MapInteger(lexicalForm string) (Integer, error) {
	vInt64, err := strconv.ParseInt(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 64)
	if err != nil {
		return Integer(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Integer(vInt64), nil
}

func (v Integer) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Integer_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (Integer) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Integer) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Integer_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
