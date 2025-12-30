package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type QuadIterator struct {
	index int
	edges statementList
}

var _ rdf.QuadIterator = &QuadIterator{}

func (i *QuadIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *QuadIterator) Err() error {
	return nil
}

func (i *QuadIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

func (i *QuadIterator) Quad() rdf.Quad {
	return i.edges[i.index].GetQuad()
}

func (i *QuadIterator) Statement() rdf.Statement {
	return i.Quad()
}

func (i *QuadIterator) StorageStatement() *Statement {
	return i.edges[i.index]
}
