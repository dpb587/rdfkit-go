package schemavalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Date struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = Date{}

func (v Date) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.Date_Class,
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
	} else if tLiteral.Datatype != schemairi.Date_Class {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
