package rdfvalue

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type LangString struct {
	Lang   string
	String string
}

var _ literalutil.CustomValue = LangString{}

func (v LangString) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    rdfiri.LangString_Datatype,
		LexicalForm: v.String,
		Tags: map[rdf.LiteralTag]string{
			rdf.LanguageLiteralTag: v.Lang,
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

	return tLiteral.Tags[rdf.LanguageLiteralTag] == v.Lang
}
