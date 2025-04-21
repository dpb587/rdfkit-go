package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewDouble(v float64) rdf.Literal {
	return xsdvalue.Double(v).AsLiteralTerm()
}

func MapDouble(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapDouble(lexicalForm)
}
