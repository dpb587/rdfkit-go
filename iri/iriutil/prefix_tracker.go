package iriutil

import (
	"maps"
	"slices"

	"github.com/dpb587/rdfkit-go/iri"
)

type UsagePrefixMapper struct {
	pm   iri.PrefixMapper
	used map[string]bool
}

var _ iri.PrefixMapper = (*UsagePrefixMapper)(nil)

func NewUsagePrefixMapper(pm iri.PrefixMapper) *UsagePrefixMapper {
	return &UsagePrefixMapper{
		pm:   pm,
		used: make(map[string]bool),
	}
}

func (p *UsagePrefixMapper) GetUsedPrefixes() []string {
	return slices.Collect(maps.Keys(p.used))
}

func (p *UsagePrefixMapper) CompactPrefix(v string) (iri.PrefixReference, bool) {
	pr, ok := p.pm.CompactPrefix(v)
	if !ok {
		return pr, false
	}

	p.used[pr.Prefix] = true

	return pr, true
}

func (p *UsagePrefixMapper) ExpandPrefix(pr iri.PrefixReference) (string, bool) {
	expanded, ok := p.pm.ExpandPrefix(pr)
	if !ok {
		return expanded, false
	}

	p.used[pr.Prefix] = true

	return expanded, true
}
