package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewGMonthDay(layout string, value time.Time) rdf.Literal {
	return xsdvalue.GMonthDay{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapGMonthDay(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapGMonthDay(lexicalForm)
}
