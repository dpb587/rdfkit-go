package iri

import (
	"maps"
	"slices"
)

type PrefixManager struct {
	ordered         PrefixMappingList
	mappingByPrefix map[string]PrefixMapping
}

var _ PrefixMapper = (*PrefixManager)(nil)

func NewPrefixManager(mappings PrefixMappingList) *PrefixManager {
	m := &PrefixManager{
		ordered:         make(PrefixMappingList, 0, len(mappings)),
		mappingByPrefix: map[string]PrefixMapping{},
	}

	m.AddPrefixMappings(mappings...)

	return m
}

func (p *PrefixManager) Clone() *PrefixManager {
	return &PrefixManager{
		ordered:         slices.Clone(p.ordered),
		mappingByPrefix: maps.Clone(p.mappingByPrefix),
	}
}

func (p *PrefixManager) AddPrefixMappings(mappings ...PrefixMapping) {
	var added int

	for _, mapping := range mappings {
		if previous, exists := p.mappingByPrefix[mapping.Prefix]; exists {
			if previous.Expanded == mapping.Expanded {
				continue
			}

			for i, existing := range p.ordered {
				if existing.Prefix == mapping.Prefix {
					p.ordered[i] = mapping
					break
				}
			}
		} else {
			p.ordered = append(p.ordered, mapping)
		}

		p.mappingByPrefix[mapping.Prefix] = mapping
		added++
	}

	if added == 0 {
		return
	}

	slices.SortFunc(p.ordered, func(a, b PrefixMapping) int {
		return len(b.Expanded) - len(a.Expanded)
	})
}

func (p *PrefixManager) DeletePrefixes(prefixes ...string) {
	var deleted int

	for _, prefix := range prefixes {
		if _, exists := p.mappingByPrefix[prefix]; !exists {
			continue
		}

		delete(p.mappingByPrefix, prefix)
		deleted++
	}

	if deleted == 0 {
		return
	}

	newOrdered := make(PrefixMappingList, 0, len(p.ordered)-deleted)

	for _, existing := range p.ordered {
		if _, exists := p.mappingByPrefix[existing.Prefix]; exists {
			newOrdered = append(newOrdered, existing)
		}
	}

	p.ordered = newOrdered
}

func (p *PrefixManager) GetPrefixMappings() PrefixMappingList {
	return slices.Clone(p.ordered)
}

func (p *PrefixManager) CompactPrefix(v string) (PrefixReference, bool) {
	for _, m := range p.ordered {
		if len(v) >= len(m.Expanded) && v[0:len(m.Expanded)] == m.Expanded {
			return PrefixReference{
				Prefix:    m.Prefix,
				Reference: v[len(m.Expanded):],
			}, true
		}
	}

	return PrefixReference{}, false
}

func (p *PrefixManager) ExpandPrefix(pr PrefixReference) (string, bool) {
	m, ok := p.mappingByPrefix[pr.Prefix]
	if !ok {
		return "", false
	}

	return m.Expanded + pr.Reference, true
}
