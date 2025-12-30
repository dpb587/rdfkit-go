package xsdtype

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type GMonth struct {
	Time   time.Time
	Layout string
}

var _ objecttypes.Value = GMonth{}

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

func (v GMonth) AsObjectValue() rdf.ObjectValue {
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
