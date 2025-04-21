package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type DatasetNodeIterator interface {
	NodeIterator

	GetGraphName() rdf.GraphNameValue
	GetTerm() rdf.Term
}
