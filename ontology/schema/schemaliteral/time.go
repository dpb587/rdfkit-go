package schemaliteral

import (
	"strings"
	"time"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
)

func CastTime(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case schemairi.Time_Class:
		return tLiteral, true
	case xsdiri.Time_Datatype:
		tValue, err := xsdvalue.MapTime(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Time_Class,
			LexicalForm: tValue.Time.Format(tValue.Layout),
		}, true
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		normalizedValue := strings.ReplaceAll(xsdutil.WhiteSpaceCollapse(tLiteral.LexicalForm), " ", "")

		for _, layout := range []struct {
			Parse  string
			Format string
		}{
			{Parse: "15:04:05.000000000Z07:00"},
			{Parse: "15:04:05.000000000Z0700"},
			{Parse: "15:04:05Z07:00"},
			{Parse: "15:04:05Z0700"},
			{Parse: "15:04:05.000000000"},
			{Parse: "15:04:05"},
			{Parse: "15:04Z07:00"},
			{Parse: "15:04Z0700"},
			{Parse: "15:04"},

			{
				Parse:  "03:04:05.0000000000PMZ07:00",
				Format: "15:04:05.0000000000Z07:00",
			},
			{
				Parse:  "03:04:05.0000000000PMZ0700",
				Format: "15:04:05.0000000000Z0700",
			},
			{
				Parse:  "03:04:05.0000000000PM",
				Format: "15:04:05.0000000000",
			},
			{
				Parse:  "03:04:05PM",
				Format: "15:04:05",
			},
			{
				Parse:  "03:04:05PMZ07:00",
				Format: "15:04:05Z07:00",
			},
			{
				Parse:  "03:04:05PMZ0700",
				Format: "15:04:05Z0700",
			},
			{
				Parse:  "03:04PM",
				Format: "15:04",
			},
			{
				Parse:  "03:04PMZ07:00",
				Format: "15:04Z07:00",
			},
			{
				Parse:  "03:04PMZ0700",
				Format: "15:04Z0700",
			},

			{
				Parse:  "3:04:05.0000000000PMZ07:00",
				Format: "15:04:05.0000000000Z07:00",
			},
			{
				Parse:  "3:04:05.0000000000PMZ0700",
				Format: "15:04:05.0000000000Z0700",
			},
			{
				Parse:  "3:04:05.0000000000PM",
				Format: "15:04:05.0000000000",
			},
			{
				Parse:  "3:04:05PM",
				Format: "15:04:05",
			},
			{
				Parse:  "3:04:05PMZ07:00",
				Format: "15:04:05Z07:00",
			},
			{
				Parse:  "3:04:05PMZ0700",
				Format: "15:04:05Z0700",
			},
			{
				Parse:  "3:04PM",
				Format: "15:04",
			},
			{
				Parse:  "3:04PMZ07:00",
				Format: "15:04Z07:00",
			},
			{
				Parse:  "3:04PMZ0700",
				Format: "15:04Z0700",
			},
		} {
			valueTime, err := time.Parse(layout.Parse, normalizedValue)
			if err == nil {
				layoutFormat := layout.Format
				if len(layoutFormat) == 0 {
					layoutFormat = layout.Parse
				}

				return rdf.Literal{
					Datatype:    schemairi.Time_Class,
					LexicalForm: valueTime.Format(layoutFormat),
				}, true
			}
		}

		return nil, false
	}

	return nil, false
}
