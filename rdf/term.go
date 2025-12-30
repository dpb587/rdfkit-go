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

func BuildTermList[S ~[]E, E Term](terms S) TermList {
	tl := make(TermList, len(terms))

	for termIdx, term := range terms {
		tl[termIdx] = term
	}

	return tl
}

//

type TermIterator interface {
	Next() bool
	Err() error
	Close() error

	Term() Term
}
