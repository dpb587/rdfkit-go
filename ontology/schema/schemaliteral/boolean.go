package schemaliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
)

func NewBoolean(v bool) rdf.ObjectValue {
	if v {
		return schemairi.True_Boolean
	}

	return schemairi.False_Boolean
}

func CastBoolean(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	if tIRI, ok := t.(rdf.IRI); ok {
		switch tIRI {
		case "http://schema.org/False", "https://schema.org/False":
			return schemairi.False_Boolean, true
		case "http://schema.org/True", "https://schema.org/True":
			return schemairi.True_Boolean, true
		}

		return nil, false
	}

	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case xsdiri.Boolean_Datatype:
		tValue, err := xsdvalue.MapBoolean(tLiteral.LexicalForm)
		if err != nil {
			return nil, false
		}

		if tValue {
			return schemairi.True_Boolean, true
		}

		return schemairi.False_Boolean, true
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		normalizedValue := xsdutil.WhiteSpaceCollapse(tLiteral.LexicalForm)

		switch normalizedValue {
		case "false",
			"False",
			"FALSE",
			"http://schema.org/False",
			"http://www.schema.org/False",
			"https://schema.org/False",
			"https://www.schema.org/False",
			"no",
			"No",
			"NO",
			"0":
			return schemairi.False_Boolean, true
		case "true",
			"True",
			"TRUE",
			"http://schema.org/True",
			"http://www.schema.org/True",
			"https://schema.org/True",
			"https://www.schema.org/True",
			"yes",
			"Yes",
			"YES",
			"1":
			return schemairi.True_Boolean, true
		}
	}

	return nil, false
}
