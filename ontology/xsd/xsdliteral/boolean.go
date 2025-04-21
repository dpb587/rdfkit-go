package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewBoolean(v bool) rdf.Literal {
	return xsdvalue.Boolean(v).AsLiteralTerm()
}

func MapBoolean(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapBoolean(lexicalForm)
}
