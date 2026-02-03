package curie

import (
	"github.com/dpb587/rdfkit-go/iri"
)

type MappingScope struct {
	Safe               bool
	DefaultPrefix      string
	DefaultPrefixEmpty bool
	Prefixes           iri.PrefixMapper
}

func (m MappingScope) CompactCURIE(v string) CURIE {
	pr, ok := m.Prefixes.CompactPrefix(v)
	if !ok {
		return CURIE{
			Safe:      m.Safe,
			Reference: string(v),
		}
	}

	c := CURIE{
		Safe:      m.Safe,
		Prefix:    pr.Prefix,
		Reference: pr.Reference,
	}

	if c.Prefix == m.DefaultPrefix && (len(c.Prefix) > 0 || m.DefaultPrefixEmpty) {
		c.Prefix = ""
		c.DefaultPrefix = true
	}

	return c
}

func (m MappingScope) ExpandCURIE(v CURIE) (string, bool) {
	if v.DefaultPrefix && (len(m.DefaultPrefix) > 0 || m.DefaultPrefixEmpty) {
		v.Prefix = m.DefaultPrefix
	}

	return m.Prefixes.ExpandPrefix(iri.PrefixReference{
		Prefix:    v.Prefix,
		Reference: v.Reference,
	})
}
