package xsdliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

func MapDuration(lexicalForm string) (literalutil.CustomValue, error) {
	return xsdvalue.MapDuration(lexicalForm)
}
