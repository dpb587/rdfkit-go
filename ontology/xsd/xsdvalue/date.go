package xsdvalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Date struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = Date{}
var _ literalutil.CustomValue = Date{}

func MapDate(lexicalForm string) (Date, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"2006-01-02",
		"2006-01-02Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return Date{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return Date{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v Date) AsTerm() rdf.Term {
	return v.AsLiteralTerm()
}

func (v Date) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Date_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (Date) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Date) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Date_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
