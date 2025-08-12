package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdfio"
)

type nodeIterator struct {
	index int
	edges nodeList
}

var _ rdfio.NodeIterator = &nodeIterator{}

func (i *nodeIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *nodeIterator) Err() error {
	return nil
}

func (i *nodeIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

// func (i *nodeIterator) GetTerm() rdf.Term {
// 	return i.edges[i.index].GetTerm()
// }

func (i *nodeIterator) GetNode() rdfio.Node {
	return i.edges[i.index]
}
