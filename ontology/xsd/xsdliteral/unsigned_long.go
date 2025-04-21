package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewUnsignedLong(v uint64) rdf.Literal {
	return xsdvalue.UnsignedLong(v).AsLiteralTerm()
}

func MapUnsignedLong(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapUnsignedLong(lexicalForm)
}
