package rdfioutil

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type graphNameIterator struct {
	index int
	edges rdf.GraphNameValueList
}

var _ rdfio.GraphNameIterator = &graphNameIterator{}

func NewStaticGraphNameIterator(edges rdf.GraphNameValueList) rdfio.GraphNameIterator {
	return &graphNameIterator{
		index: -1,
		edges: edges,
	}
}

func (i *graphNameIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *graphNameIterator) Err() error {
	return nil
}

func (i *graphNameIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

func (i *graphNameIterator) GetGraphName() rdf.GraphNameValue {
	return i.edges[i.index]
}
