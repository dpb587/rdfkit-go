package rdf

type TermKind int

const (
	TermKindBlankNode TermKind = iota
	TermKindIRI
	TermKindLiteral
)

//

// Term represents a value that may be used in some position of a statement.
//
// This is a closed interface. See [BlankNode], [IRI], and [Literal].
type Term interface {
	TermKind() TermKind
	TermEquals(a Term) bool

	isTermBuiltin()
}

//

type TermMatcher interface {
	MatchTerm(t Term) bool
}

//

type TermMatcherFunc func(t Term) bool

var _ TermMatcher = TermMatcherFunc(nil)

func (f TermMatcherFunc) MatchTerm(t Term) bool {
	return f(t)
}

//

type TermList []Term

//

type TermIterator interface {
	Next() bool
	Err() error
	Close() error

	Term() Term
}
