package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewTime(layout string, value time.Time) rdf.Literal {
	return xsdvalue.Time{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapTime(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapTime(lexicalForm)
}
