package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type GDay struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = GDay{}
var _ literalutil.CustomValue = GDay{}

func MapGDay(lexicalForm string) (GDay, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"02",
		"02Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return GDay{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return GDay{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v GDay) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v GDay) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.GDay_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (GDay) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v GDay) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.GDay_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
