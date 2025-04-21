package literalutil

import "github.com/dpb587/rdfkit-go/rdf"

type CustomValue interface {
	// previously from rdf.TermBase; probably resurrect?
	interface {
		TermKind() rdf.TermKind
		TermEquals(a rdf.Term) bool
	}

	AsLiteralTerm() rdf.Literal
}
