package terms

import "github.com/dpb587/rdfkit-go/rdf"

type iterator struct {
	terms rdf.TermList
	idx   int
}

var _ rdf.TermIterator = &iterator{}

func NewIterator(terms rdf.TermList) rdf.TermIterator {
	return &iterator{
		terms: terms,
		idx:   -1,
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

	return it.idx < len(it.terms)
}

func (it *iterator) Term() rdf.Term {
	return it.terms[it.idx]
}
