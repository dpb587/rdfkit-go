package schemaliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func NewText(v string) rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.Text_Class,
		LexicalForm: v,
	}
}

func CastText(t rdf.Term, _ CastOptions) (rdf.ObjectValue, bool) {
	switch t := t.(type) {
	case rdf.Literal:
		switch t.Datatype {
		case schemairi.Text_Class:
			return t, true
		case rdfiri.LangString_Datatype,
			xsdiri.DateTimeStamp_Datatype,
			xsdiri.DateTime_Datatype,
			xsdiri.Date_Datatype,
			xsdiri.Decimal_Datatype,
			xsdiri.Integer_Datatype,
			xsdiri.String_Datatype:
			return rdf.Literal{
				Datatype:    schemairi.Text_Class,
				LexicalForm: t.LexicalForm,
			}, true
		}
	case rdf.IRI:
		return rdf.Literal{
			Datatype:    schemairi.Text_Class,
			LexicalForm: string(t),
		}, true
	}

	return nil, false
}
