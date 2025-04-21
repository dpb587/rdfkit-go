package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewFloat(v float32) rdf.Literal {
	return xsdvalue.Float(v).AsLiteralTerm()
}

func MapFloat(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapFloat(lexicalForm)
}
