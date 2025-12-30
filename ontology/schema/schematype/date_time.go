package schematype

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type DateTime struct {
	Time   time.Time
	Layout string
}

var _ objecttypes.Value = DateTime{}

func (v DateTime) AsObjectValue() rdf.ObjectValue {
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
