package rdfvalue

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type HTML string

var _ literalutil.CustomValue = HTML("")

func MapHTML(s HTML) (HTML, error) {
	return HTML(s), nil
}

func (v HTML) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    rdfiri.HTML_Datatype,
		LexicalForm: string(v),
	}
}

func (HTML) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v HTML) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != rdfiri.HTML_Datatype {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
