package xsdliteral

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func NewDateTimeStamp(layout string, value time.Time) rdf.Literal {
	return xsdvalue.DateTimeStamp{
		Time:   value,
		Layout: layout,
	}.AsLiteralTerm()
}

func MapDateTimeStamp(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapDateTimeStamp(lexicalForm)
}
