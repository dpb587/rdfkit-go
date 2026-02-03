package rdf

func BuildTermList[S ~[]E, E Term](terms S) TermList {
	tl := make(TermList, len(terms))

	for termIdx, term := range terms {
		tl[termIdx] = term
	}

	return tl
}
