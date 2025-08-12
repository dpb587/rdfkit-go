package termutil

import "github.com/dpb587/rdfkit-go/rdf"

type CustomValue interface {
	TermKind() rdf.TermKind
	TermEquals(a rdf.Term) bool
	AsTerm() rdf.Term
}
