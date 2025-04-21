package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewByte(v byte) rdf.Literal {
	return xsdvalue.Byte(v).AsLiteralTerm()
}

func MapByte(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapByte(lexicalForm)
}
