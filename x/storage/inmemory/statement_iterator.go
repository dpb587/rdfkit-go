package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type statementIterator struct {
	index int
	edges statementList
}

var _ rdfio.StatementIterator = &statementIterator{}
var _ rdfio.GraphStatementIterator = &statementIterator{}
var _ rdfio.DatasetStatementIterator = &statementIterator{}

func (i *statementIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *statementIterator) Err() error {
	return nil
}

func (i *statementIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

func (i *statementIterator) GetGraphName() rdf.GraphNameValue {
	return i.edges[i.index].g.t.(rdf.GraphNameValue)
}

func (i *statementIterator) GetTriple() rdf.Triple {
	return i.edges[i.index].GetTriple()
}

func (i *statementIterator) GetStatement() rdfio.Statement {
	return i.edges[i.index]
}
