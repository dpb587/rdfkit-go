package iriutil

import "github.com/dpb587/rdfkit-go/rdf"

type PrefixManager interface {
	CompactPrefix(v rdf.IRI) (string, string, bool)
	ExpandPrefix(prefix, reference string) (rdf.IRI, bool)
}
