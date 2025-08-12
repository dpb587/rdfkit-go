package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Time struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = Time{}
var _ literalutil.CustomValue = Time{}

func MapTime(lexicalForm string) (Time, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"15:04:05",
		"15:04:05.000000000",
		"15:04:05Z",
		"15:04:05.000000000Z",
		"15:04:05Z07:00",
		"15:04:05.000000000Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return Time{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return Time{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v Time) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v Time) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Time_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (Time) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Time) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Time_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
