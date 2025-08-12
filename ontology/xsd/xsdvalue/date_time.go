package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type DateTime struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = DateTime{}
var _ literalutil.CustomValue = DateTime{}

func MapDateTime(lexicalForm string) (DateTime, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.000000000",
		"2006-01-02T15:04:05.000000000Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return DateTime{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return DateTime{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v DateTime) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v DateTime) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.DateTime_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (DateTime) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v DateTime) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.DateTime_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
