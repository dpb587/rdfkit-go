package inmemory

import "github.com/dpb587/rdfkit-go/rdf"

type termIterator struct {
	index int
	edges nodeList
}

var _ rdf.TermIterator = &termIterator{}

func (i *termIterator) Close() error {
	i.edges = nil

	return nil
}

func (i *termIterator) Err() error {
	return nil
}

func (i *termIterator) Next() bool {
	if i.index >= len(i.edges)-1 {
		return false
	}

	i.index++

	return true
}

func (i *termIterator) Term() rdf.Term {
	return i.edges[i.index].t
}
