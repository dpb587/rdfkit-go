package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewInteger(v int64) rdf.Literal {
	return xsdvalue.Integer(v).AsLiteralTerm()
}

func MapInteger(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapInteger(lexicalForm)
}
