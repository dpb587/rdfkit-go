package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewString(v string) rdf.Literal {
	return xsdvalue.String(v).AsLiteralTerm()
}

func MapString(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapString(lexicalForm)
}
