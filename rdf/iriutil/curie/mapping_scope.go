package curie

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type MappingScope struct {
	Safe               bool
	DefaultPrefix      string
	DefaultPrefixEmpty bool
	Prefixes           iriutil.PrefixManager
}

func (m MappingScope) CompactCURIE(v rdf.IRI) CURIE {
	prefix, reference, ok := m.Prefixes.CompactPrefix(v)
	if !ok {
		return CURIE{
			Safe:      m.Safe,
			Reference: string(v),
		}
	}

	c := CURIE{
		Safe:      m.Safe,
		Prefix:    prefix,
		Reference: reference,
	}

	if c.Prefix == m.DefaultPrefix && (len(c.Prefix) > 0 || m.DefaultPrefixEmpty) {
		c.Prefix = ""
		c.DefaultPrefix = true
	}

	return c
}

func (m MappingScope) ExpandCURIE(v CURIE) (rdf.IRI, bool) {
	if v.DefaultPrefix && (len(m.DefaultPrefix) > 0 || m.DefaultPrefixEmpty) {
		v.Prefix = m.DefaultPrefix
	}

	return m.Prefixes.ExpandPrefix(v.Prefix, v.Reference)
}
