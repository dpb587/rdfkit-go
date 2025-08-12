package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type GMonth struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = GMonth{}
var _ literalutil.CustomValue = GMonth{}

func MapGMonth(lexicalForm string) (GMonth, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"01",
		"01Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return GMonth{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return GMonth{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v GMonth) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v GMonth) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.GMonth_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (GMonth) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v GMonth) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.GMonth_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
