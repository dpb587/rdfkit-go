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

type Double float64

var _ termutil.CustomValue = Double(0)
var _ literalutil.CustomValue = Double(0)

func MapDouble(lexicalForm string) (Double, error) {
	vFloat64, err := strconv.ParseFloat(xsdutil.WhiteSpaceCollapse(lexicalForm), 64)
	if err != nil {
		return Double(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Double(vFloat64), nil
}

func (v Double) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v Double) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Double_Datatype,
		LexicalForm: strconv.FormatFloat(float64(v), 'f', -1, 64),
	}
}

func (Double) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Double) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Double_Datatype {
		return false
	}

	return strconv.FormatFloat(float64(v), 'f', -1, 64) == tLiteral.LexicalForm
}
