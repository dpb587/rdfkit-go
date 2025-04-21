package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type DatasetStatementIterator interface {
	StatementIterator

	GetGraphName() rdf.GraphNameValue
	GetTriple() rdf.Triple
}
