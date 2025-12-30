package terms

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type iteratorIterator struct {
	iters []rdf.TermIterator
}

var _ rdf.TermIterator = &iteratorIterator{}

func NewIteratorIterator(iters ...rdf.TermIterator) rdf.TermIterator {
	return &iteratorIterator{
		iters: iters,
	}
}

func (i *iteratorIterator) Close() error {
	for _, iter := range i.iters {
		iter.Close()
	}

	return nil
}

func (i *iteratorIterator) Err() error {
	if len(i.iters) > 0 {
		return i.iters[0].Err()
	}

	return nil
}

func (i *iteratorIterator) Next() bool {
	for {
		if len(i.iters) == 0 {
			return false
		} else if i.iters[0].Err() != nil {
			return false
		} else if i.iters[0].Next() {
			return true
		}

		i.iters = i.iters[1:]
	}
}

func (i *iteratorIterator) Term() rdf.Term {
	return i.iters[0].Term()
}
