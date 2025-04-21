package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewDecimal(v float64) rdf.Literal {
	return xsdvalue.Decimal(v).AsLiteralTerm()
}

func MapDecimal(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapDecimal(lexicalForm)
}
