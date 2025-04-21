package rdfliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfvalue"
	"github.com/dpb587/rdfkit-go/rdf"
)

func NewHTML(v string) rdf.Literal {
	return rdfvalue.HTML(v).AsLiteralTerm()
}
