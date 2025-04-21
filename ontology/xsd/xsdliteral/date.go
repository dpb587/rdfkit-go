package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewDate(layout string, value time.Time) rdf.Literal {
	return xsdvalue.Date{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapDate(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapDate(lexicalForm)
}
