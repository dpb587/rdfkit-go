package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewUnsignedInt(v uint32) rdf.Literal {
	return xsdvalue.UnsignedInt(v).AsLiteralTerm()
}

func MapUnsignedInt(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapUnsignedInt(lexicalForm)
}
