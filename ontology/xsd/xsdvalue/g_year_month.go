package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type GYearMonth struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = GYearMonth{}
var _ literalutil.CustomValue = GYearMonth{}

func MapGYearMonth(lexicalForm string) (GYearMonth, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"2006-01",
		"2006-01Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return GYearMonth{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return GYearMonth{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v GYearMonth) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v GYearMonth) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.GYearMonth_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (GYearMonth) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v GYearMonth) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.GYearMonth_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
