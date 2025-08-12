package schemavalue

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type DateTime struct {
	Time   time.Time
	Layout string
}

var _ termutil.CustomValue = DateTime{}

func (v DateTime) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.DateTime_Class,
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
	} else if tLiteral.Datatype != schemairi.DateTime_Class {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
