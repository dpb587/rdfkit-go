package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewHexBinary(v []byte) rdf.Literal {
	return xsdvalue.HexBinary(v).AsLiteralTerm()
}

func MapHexBinary(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapHexBinary(lexicalForm)
}
