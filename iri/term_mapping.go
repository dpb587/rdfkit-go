package iri

import (
	"strings"
)

type TermMapping struct {
	Term     string
	Expanded string
}

//

type TermMappingList []TermMapping

//

func CompareTermMappingByTerm(a, b TermMapping) int {
	return strings.Compare(a.Term, b.Term)
}
