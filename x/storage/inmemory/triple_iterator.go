package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type TripleIterator struct {
	index int
	edges statementList
}

var _ rdf.TripleIterator = &TripleIterator{}

func (i *TripleIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *TripleIterator) Err() error {
	return nil
}

func (i *TripleIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

func (i *TripleIterator) Triple() rdf.Triple {
	return i.edges[i.index].GetQuad().Triple
}

func (i *TripleIterator) Statement() rdf.Statement {
	return i.Triple()
}

func (i *TripleIterator) StorageStatement() *Statement {
	return i.edges[i.index]
}
