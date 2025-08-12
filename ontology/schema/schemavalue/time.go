package schemavalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Time struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = Time{}

func (v Time) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.Time_Class,
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
	} else if tLiteral.Datatype != schemairi.Time_Class {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
