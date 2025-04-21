package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewAnyURI(v string) rdf.Literal {
	return xsdvalue.AnyURI(v).AsLiteralTerm()
}

func MapAnyURI(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapAnyURI(lexicalForm)
}
