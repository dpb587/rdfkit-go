package triples

import "github.com/dpb587/rdfkit-go/rdf"

type iterator struct {
	triples rdf.TripleList
	idx     int
}

var _ rdf.TripleIterator = &iterator{}

func NewIterator(triples rdf.TripleList) rdf.TripleIterator {
	return &iterator{
		triples: triples,
		idx:     -1,
	}
}

func (it *iterator) Close() error {
	return nil
}

func (it *iterator) Err() error {
	return nil
}

func (it *iterator) Next() bool {
	it.idx++

	return it.idx < len(it.triples)
}

func (it *iterator) Triple() rdf.Triple {
	return it.triples[it.idx]
}

func (it *iterator) Statement() rdf.Statement {
	return it.Triple()
}
