package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewDateTime(layout string, value time.Time) rdf.Literal {
	return xsdvalue.DateTime{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapDateTime(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapDateTime(lexicalForm)
}
