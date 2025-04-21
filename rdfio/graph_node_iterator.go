package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type GraphNodeIterator interface {
	NodeIterator

	GetTerm() rdf.Term
}
