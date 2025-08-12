package schemavalue

import (
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type URL rdf.IRI

var _ termutil.CustomValue = URL("")

func (v URL) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.URL_Class,
		LexicalForm: string(v),
	}
}

func (URL) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v URL) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.URL_Class {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
