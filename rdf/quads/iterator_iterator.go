package quads

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type iteratorIterator struct {
	iters []rdf.QuadIterator
	err   error
}

var _ rdf.QuadIterator = &iteratorIterator{}

func NewIteratorIterator(iters ...rdf.QuadIterator) rdf.QuadIterator {
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
	return i.err
}

func (i *iteratorIterator) Next() bool {
	for {
		if len(i.iters) == 0 {
			return false
		} else if i.iters[0].Next() {
			return true
		} else if v := i.iters[0].Err(); v != nil {
			i.err = v

			return false
		}

		i.iters = i.iters[1:]
	}
}

func (i *iteratorIterator) Quad() rdf.Quad {
	return i.iters[0].Quad()
}

func (i *iteratorIterator) Statement() rdf.Statement {
	return i.Quad()
}
