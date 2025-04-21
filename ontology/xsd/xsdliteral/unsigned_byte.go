package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewUnsignedByte(v uint8) rdf.Literal {
	return xsdvalue.UnsignedByte(v).AsLiteralTerm()
}

func MapUnsignedByte(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapUnsignedByte(lexicalForm)
}
