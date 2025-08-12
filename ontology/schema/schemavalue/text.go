package schemavalue

import (
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Text string

var _ termutil.CustomValue = Text("")

func (v Text) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.Text_Class,
		LexicalForm: string(v),
	}
}

func (Text) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Text) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.Text_Class {
		return false
	}

	return tLiteral.LexicalForm == string(v)
}
