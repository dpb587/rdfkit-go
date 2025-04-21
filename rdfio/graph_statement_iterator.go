package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type GraphStatementIterator interface {
	StatementIterator

	GetTriple() rdf.Triple
}
