package quads

import "github.com/dpb587/rdfkit-go/rdf"

type iterator struct {
	quads rdf.QuadList
	idx   int
}

var _ rdf.QuadIterator = &iterator{}

func NewIterator(quads rdf.QuadList) rdf.QuadIterator {
	return &iterator{
		quads: quads,
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

	return it.idx < len(it.quads)
}

func (it *iterator) Quad() rdf.Quad {
	return it.quads[it.idx]
}

func (it *iterator) Statement() rdf.Statement {
	return it.Quad()
}
