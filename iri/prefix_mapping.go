package iri

import "strings"

type PrefixMapping struct {
	// Prefix is the short name which should match XML's QName conventions and exclude a trailing colon.
	Prefix string

	// Expanded is the equivalent IRI which should end with a slash or number sign.
	Expanded string
}

//

type PrefixMappingList []PrefixMapping

//

func ComparePrefixMappingByPrefix(a, b PrefixMapping) int {
	return strings.Compare(a.Prefix, b.Prefix)
}
