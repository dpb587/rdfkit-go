package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewGDay(layout string, value time.Time) rdf.Literal {
	return xsdvalue.GDay{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapGDay(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapGDay(lexicalForm)
}
