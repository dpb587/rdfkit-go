package iriutil

import (
	"strings"
)

type PrefixMapping struct {
	Prefix   string
	Expanded string
}

type PrefixMappingList []PrefixMapping

//

func ComparePrefixMappingByPrefix(a, b PrefixMapping) int {
	return strings.Compare(a.Prefix, b.Prefix)
}
