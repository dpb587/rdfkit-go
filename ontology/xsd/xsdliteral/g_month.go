package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewGMonth(layout string, value time.Time) rdf.Literal {
	return xsdvalue.GMonth{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapGMonth(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapGMonth(lexicalForm)
}
