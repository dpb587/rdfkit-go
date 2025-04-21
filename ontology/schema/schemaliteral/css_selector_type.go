package schemaliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func CastCssSelectorType(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return nil, false
	}

	switch tLiteral.Datatype {
	case rdfiri.LangString_Datatype,
		xsdiri.String_Datatype:
		return rdf.Literal{
			Datatype:    schemairi.CssSelectorType_Class,
			LexicalForm: tLiteral.LexicalForm,
		}, true
	}

	return nil, false
}
