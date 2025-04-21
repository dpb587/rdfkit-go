package schemaliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
)

func CastDateTime(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case schemairi.DateTime_Class:
		return tLiteral, true
	case xsdiri.DateTimeStamp_Datatype:
		tValue, err := xsdvalue.MapDateTimeStamp(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.DateTime_Class,
			LexicalForm: tValue.Time.Format(tValue.Layout),
		}, true
	case xsdiri.DateTime_Datatype:
		tValue, err := xsdvalue.MapDateTime(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.DateTime_Class,
			LexicalForm: tValue.Time.Format(tValue.Layout),
		}, true
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		normalizedValue := xsdutil.WhiteSpaceCollapse(tLiteral.LexicalForm)

		for _, layout := range []string{
			"2006-01-02T15:04:05.000000000Z07:00",
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05Z0700",
			"2006-01-02T15:04:05.000000000",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04Z07:00",
			"2006-01-02T15:04Z0700",
			"2006-01-02T15:04",
			"2006-01-02T15Z07:00",
			"2006-01-02T15Z0700",
			"2006-01-02T15",
		} {
			valueTime, err := time.Parse(layout, normalizedValue)
			if err == nil {
				return rdf.Literal{
					Datatype:    schemairi.DateTime_Class,
					LexicalForm: valueTime.Format(layout),
				}, true
			}
		}
	}

	return nil, false
}
