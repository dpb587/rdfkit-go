package iriutil

import (
	"maps"

	"github.com/dpb587/rdfkit-go/rdf"
)

type PrefixMap map[string]rdf.IRI

var _ PrefixManager = PrefixMap{}

func NewPrefixMap(entries ...PrefixMapping) PrefixMap {
	v := PrefixMap{}

	for _, i := range entries {
		v[i.Prefix] = i.Expanded
	}

	return v
}

func (p PrefixMap) NewPrefixMap(entries ...PrefixMapping) PrefixMap {
	next := maps.Clone(p)

	for _, i := range entries {
		next[i.Prefix] = i.Expanded
	}

	return next
}

func (p PrefixMap) AsPrefixMappingList() PrefixMappingList {
	var res PrefixMappingList

	for prefix, expanded := range p {
		res = append(res, PrefixMapping{
			Prefix:   prefix,
			Expanded: expanded,
		})
	}

	return res
}

func (p PrefixMap) CompactPrefix(v rdf.IRI) (string, string, bool) {
	var matchPrefix string
	var matchLen int

	for prefix, expanded := range p {
		if len(v) < len(expanded) || v[:len(expanded)] != expanded {
			continue
		} else if len(expanded) > matchLen {
			matchPrefix = prefix
			matchLen = len(expanded)
		}
	}

	if matchLen == 0 {
		return "", "", false
	}

	return matchPrefix, string(v)[matchLen:], true
}

func (p PrefixMap) ExpandPrefix(prefix, reference string) (rdf.IRI, bool) {
	expanded, ok := p[prefix]
	if !ok {
		return "", false
	}

	return expanded + rdf.IRI(reference), true
}
