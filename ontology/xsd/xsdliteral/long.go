package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewLong(v int64) rdf.Literal {
	return xsdvalue.Long(v).AsLiteralTerm()
}

func MapLong(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapLong(lexicalForm)
}
