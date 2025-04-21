package rdf

type TermKind int

const (
	TermKindBlankNode TermKind = iota
	TermKindIRI
	TermKindLiteral
)

//

type Term interface {
	TermKind() TermKind
	TermEquals(a Term) bool

	isTermBuiltin()
}

type TermList []Term
