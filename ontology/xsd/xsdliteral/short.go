package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewShort(v int16) rdf.Literal {
	return xsdvalue.Short(v).AsLiteralTerm()
}

func MapShort(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapShort(lexicalForm)
}
