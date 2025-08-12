package schemavalue

import (
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Boolean bool

var _ termutil.CustomValue = Boolean(false)

func (v Boolean) AsTerm() rdf.Term {
	if v {
		return schemairi.True_Boolean
	}

	return schemairi.False_Boolean
}

func (Boolean) TermKind() rdf.TermKind {
	return rdf.TermKindIRI
}

func (v Boolean) TermEquals(t rdf.Term) bool {
	tIRI, ok := t.(rdf.IRI)
	if !ok {
		return false
	} else if tIRI == schemairi.True_Boolean {
		return bool(v)
	} else if tIRI == schemairi.False_Boolean {
		return !bool(v)
	}

	return false
}
