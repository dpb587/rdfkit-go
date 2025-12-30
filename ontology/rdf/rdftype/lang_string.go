package rdftype

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type LangString struct {
	Lang   string
	String string
}

var _ objecttypes.Value = LangString{}

func (v LangString) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    rdfiri.LangString_Datatype,
		LexicalForm: v.String,
		Tag: rdf.LanguageLiteralTag{
			Language: v.Lang,
		},
	}
}

func (LangString) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v LangString) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != rdfiri.LangString_Datatype {
		return false
	} else if tLiteral.LexicalForm != v.String {
		return false
	}

	if langTag, ok := tLiteral.Tag.(rdf.LanguageLiteralTag); ok {
		return langTag.Language == v.Lang
	}

	return false
}
