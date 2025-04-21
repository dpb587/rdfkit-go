package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewInt(v int32) rdf.Literal {
	return xsdvalue.Int(v).AsLiteralTerm()
}

func MapInt(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapInt(lexicalForm)
}
