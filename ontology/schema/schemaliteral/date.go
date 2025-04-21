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

func CastDate(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case schemairi.Date_Class:
		return tLiteral, true
	case xsdiri.Date_Datatype:
		tValue, err := xsdvalue.MapDate(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Date_Class,
			LexicalForm: tValue.Time.Format(tValue.Layout),
		}, true
	case xsdiri.GYearMonth_Datatype:
		tValue, err := xsdvalue.MapGYearMonth(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Date_Class,
			LexicalForm: tValue.Time.Format(tValue.Layout),
		}, true
	case xsdiri.GYear_Datatype:
		tValue, err := xsdvalue.MapGYear(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Date_Class,
			LexicalForm: tValue.Time.Format(tValue.Layout),
		}, true
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		normalizedValue := xsdutil.WhiteSpaceCollapse(tLiteral.LexicalForm)

		for _, layout := range []string{
			"2006-01-02Z07:00",
			"2006-01-02",
			"2006-01Z07:00",
			"2006-01",
			"2006Z07:00",
			"2006",
		} {
			parsed, err := time.Parse(layout, normalizedValue)
			if err != nil {
				continue
			}

			return rdf.Literal{
				Datatype:    schemairi.Date_Class,
				LexicalForm: parsed.Format(layout),
			}, true
		}
	}

	return nil, false
}
