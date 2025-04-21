package rdfliteral

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfvalue"
	"github.com/dpb587/rdfkit-go/rdf"
)

func NewLangString(lang, v string) rdf.Literal {
	return rdfvalue.LangString{
		Lang:   lang,
		String: v,
	}.AsLiteralTerm()
}
