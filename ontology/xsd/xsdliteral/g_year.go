package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewGYear(layout string, value time.Time) rdf.Literal {
	return xsdvalue.GYear{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapGYear(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapGYear(lexicalForm)
}
