package iriutil

import (
	"strings"

	"github.com/dpb587/rdfkit-go/rdf"
)

type TermMapping struct {
	Term     string
	Expanded rdf.IRI
}

type TermMappingList []TermMapping

//

func CompareTermMappingByTerm(a, b TermMapping) int {
	return strings.Compare(a.Term, b.Term)
}
