package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewGYearMonth(layout string, value time.Time) rdf.Literal {
	return xsdvalue.GYearMonth{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapGYearMonth(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapGYearMonth(lexicalForm)
}
