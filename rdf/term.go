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

type TermList []Term

//

func BuildTermList[S ~[]E, E Term](terms S) TermList {
	tl := make(TermList, len(terms))

	for termIdx, term := range terms {
		tl[termIdx] = term
	}

	return tl
}
