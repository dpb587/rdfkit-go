package xsdvalue

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type Float float32

var _ literalutil.CustomValue = Float(0)

func MapFloat(lexicalForm string) (Float, error) {
	vFloat64, err := strconv.ParseFloat(xsdutil.WhiteSpaceCollapse(lexicalForm), 32)
	if err != nil {
		return Float(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Float(vFloat64), nil
}

func (v Float) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Float_Datatype,
		LexicalForm: strconv.FormatFloat(float64(v), 'f', -1, 32),
	}
}

func (Float) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Float) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Float_Datatype {
		return false
	}

	return strconv.FormatFloat(float64(v), 'f', -1, 32) == tLiteral.LexicalForm
}
