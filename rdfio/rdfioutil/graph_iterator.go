package rdfioutil

import "github.com/dpb587/rdfkit-go/rdfio"

type staticGraphIterator struct {
	index  int
	graphs []rdfio.Graph
}

var _ rdfio.GraphIterator = &staticGraphIterator{}

func NewStaticGraphIterator(graphs []rdfio.Graph) rdfio.GraphIterator {
	return &staticGraphIterator{
		index:  -1,
		graphs: graphs,
	}
}

func (i *staticGraphIterator) Close() error {
	i.graphs = nil

	return nil
}

func (i *staticGraphIterator) Err() error {
	return nil
}

func (i *staticGraphIterator) Next() bool {
	if i.index >= len(i.graphs)-1 {
		return false
	}

	i.index++

	return true
}

func (i *staticGraphIterator) GetGraph() rdfio.Graph {
	return i.graphs[i.index]
}
