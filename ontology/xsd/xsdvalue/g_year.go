package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type GYear struct {
	Time   time.Time
	Layout string
}

var _ literalutil.CustomValue = GYear{}

func MapGYear(lexicalForm string) (GYear, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"2006",
		"2006Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return GYear{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return GYear{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v GYear) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.GYear_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (GYear) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v GYear) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.GYear_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
