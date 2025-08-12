package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type GMonthDay struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = GMonthDay{}
var _ literalutil.CustomValue = GMonthDay{}

func MapGMonthDay(lexicalForm string) (GMonthDay, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"01-02",
		"01-02Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return GMonthDay{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return GMonthDay{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v GMonthDay) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v GMonthDay) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.GMonthDay_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (GMonthDay) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v GMonthDay) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.GMonthDay_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
