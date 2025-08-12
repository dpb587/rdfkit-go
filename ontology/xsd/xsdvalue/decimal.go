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

type Decimal float64

var _ termutil.CustomValue = Decimal(0)
var _ literalutil.CustomValue = Decimal(0)

func MapDecimal(lexicalForm string) (Decimal, error) {
	vFloat64, err := strconv.ParseFloat(xsdutil.WhiteSpaceCollapse(lexicalForm), 64)
	if err != nil {
		return Decimal(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Decimal(vFloat64), nil
}

func (v Decimal) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v Decimal) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Decimal_Datatype,
		LexicalForm: strconv.FormatFloat(float64(v), 'f', -1, 64),
	}
}

func (Decimal) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Decimal) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Decimal_Datatype {
		return false
	}

	return strconv.FormatFloat(float64(v), 'f', -1, 64) == tLiteral.LexicalForm
}
