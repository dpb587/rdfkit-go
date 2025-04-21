package iriutil

import "github.com/dpb587/rdfkit-go/rdf"

type PrefixManager interface {
	CompactIRI(v rdf.IRI) (string, string, bool)
	ExpandIRI(prefix, reference string) (rdf.IRI, bool)
}
