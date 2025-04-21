package schemaliteral

import (
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
)

func NewFloat(v float64) rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.Float_Class,
		LexicalForm: strconv.FormatFloat(v, 'f', -1, 64),
	}
}

func CastFloat(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case schemairi.Float_Class:
		return tLiteral, true
	case xsdiri.Decimal_Datatype:
		tValue, err := xsdvalue.MapDecimal(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Float_Class,
			LexicalForm: strconv.FormatFloat(float64(tValue), 'f', -1, 64),
		}, true
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		normalizedValue := xsdutil.WhiteSpaceCollapse(tLiteral.LexicalForm)

		tFloat64, err := strconv.ParseFloat(normalizedValue, 64)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Float_Class,
			LexicalForm: strconv.FormatFloat(tFloat64, 'f', -1, 64),
		}, true
	}

	return nil, false
}
