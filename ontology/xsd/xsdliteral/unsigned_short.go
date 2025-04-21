package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewUnsignedShort(v uint16) rdf.Literal {
	return xsdvalue.UnsignedShort(v).AsLiteralTerm()
}

func MapUnsignedShort(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapUnsignedShort(lexicalForm)
}
