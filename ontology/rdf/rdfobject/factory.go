package rdfobject

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdftype"
	"github.com/dpb587/rdfkit-go/rdf"
)

func HTML(v string) rdf.ObjectValue {
	return rdftype.HTML(v).AsObjectValue()
}

func LangString(lang, v string) rdf.ObjectValue {
	return rdftype.LangString{
		Lang:   lang,
		String: v,
	}.AsObjectValue()
}
