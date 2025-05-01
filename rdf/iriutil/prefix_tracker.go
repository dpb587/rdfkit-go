package iriutil

import "github.com/dpb587/rdfkit-go/rdf"

type PrefixTracker struct {
	pm      PrefixManager
	tracked map[string]bool
}

var _ PrefixManager = (*PrefixTracker)(nil)

func NewPrefixTracker(pm PrefixManager) *PrefixTracker {
	return &PrefixTracker{
		pm:      pm,
		tracked: make(map[string]bool),
	}
}

func (p *PrefixTracker) GetUsedPrefixMappings() PrefixMappingList {
	var res PrefixMappingList

	for prefix := range p.tracked {
		expanded, ok := p.pm.ExpandPrefix(prefix, "")
		if !ok {
			continue // weird
		}

		res = append(res, PrefixMapping{
			Prefix:   prefix,
			Expanded: string(expanded),
		})
	}

	return res
}

func (p *PrefixTracker) CompactPrefix(v rdf.IRI) (string, string, bool) {
	prefix, reference, ok := p.pm.CompactPrefix(v)
	if !ok {
		return "", "", false
	}

	p.tracked[prefix] = true

	return prefix, reference, true
}

func (p *PrefixTracker) ExpandPrefix(prefix, reference string) (rdf.IRI, bool) {
	expanded, ok := p.pm.ExpandPrefix(prefix, reference)
	if !ok {
		return "", false
	}

	p.tracked[prefix] = true

	return expanded, true
}
