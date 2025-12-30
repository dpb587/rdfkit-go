package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type StatementIterator struct {
	index int
	edges statementList
}

var _ rdf.QuadIterator = &StatementIterator{}

func (i *StatementIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *StatementIterator) Err() error {
	return nil
}

func (i *StatementIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

func (i *StatementIterator) Quad() rdf.Quad {
	return i.edges[i.index].GetQuad()
}

func (i *StatementIterator) Statement() rdf.Statement {
	return i.Quad()
}

func (i *StatementIterator) GetStatement() *Statement {
	return i.edges[i.index]
}
