package xsdtype

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type DateTimeStamp struct {
	Time   time.Time
	Layout string
}

var _ objecttypes.Value = DateTimeStamp{}

func MapDateTimeStamp(lexicalForm string) (DateTimeStamp, error) {
	lexicalForm = xsdutil.WhiteSpaceCollapse(lexicalForm)

	for _, layout := range []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.000000000Z07:00",
	} {
		parsed, err := time.Parse(layout, lexicalForm)
		if err == nil {
			return DateTimeStamp{
				Time:   parsed,
				Layout: layout,
			}, nil
		}
	}

	return DateTimeStamp{}, rdf.ErrLiteralLexicalFormNotValid
}

func (v DateTimeStamp) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    xsdiri.DateTimeStamp_Datatype,
		LexicalForm: v.Time.Format(v.Layout),
	}
}

func (DateTimeStamp) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v DateTimeStamp) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.DateTimeStamp_Datatype {
		return false
	}

	return v.Time.Format(v.Layout) == tLiteral.LexicalForm
}
