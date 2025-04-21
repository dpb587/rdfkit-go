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

func NewInteger(v int64) rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.Integer_Class,
		LexicalForm: strconv.FormatInt(v, 10),
	}
}

func CastInteger(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case schemairi.Integer_Class:
		return tLiteral, true
	case xsdiri.Integer_Datatype:
		tValue, err := xsdvalue.MapInteger(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Integer_Class,
			LexicalForm: strconv.FormatInt(int64(tValue), 10),
		}, true
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		normalizedValue := xsdutil.WhiteSpaceCollapse(tLiteral.LexicalForm)

		tInt64, err := strconv.ParseInt(normalizedValue, 0, 64)
		if err != nil {
			return nil, false
		}

		return rdf.Literal{
			Datatype:    schemairi.Integer_Class,
			LexicalForm: strconv.FormatInt(tInt64, 10),
		}, true
	}

	return nil, false
}
