package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewBase64Binary(v []byte) rdf.Literal {
	return xsdvalue.Base64Binary(v).AsLiteralTerm()
}

func MapBase64Binary(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapBase64Binary(lexicalForm)
}
