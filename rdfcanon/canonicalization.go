package rdfcanon

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type canonicalizedQuad struct {
	originalIndex int64
	canonical     *rdf.Quad
	encoded       []byte
}

//

type Canonicalization struct {
	blankNodeStringer blanknodeutil.Stringer
	hasCanonicalQuad  bool
	nquads            []canonicalizedQuad
}

func (c *Canonicalization) GetBlankNodeIdentifier(bn rdf.BlankNode) string {
	return c.blankNodeStringer.GetBlankNodeIdentifier(bn)
}

func (c *Canonicalization) AsQuads() rdf.QuadList {
	if !c.hasCanonicalQuad {
		panic("canonical quads was not enabled")
	}

	quads := make(rdf.QuadList, 0, len(c.nquads))

	for _, nquad := range c.nquads {
		quads = append(quads, *nquad.canonical)
	}

	return quads
}

func (c *Canonicalization) NewIterator() *CanonicalizedIterator {
	return &CanonicalizedIterator{
		c:   c,
		idx: -1,
	}
}

func (c *Canonicalization) NewQuadIterator() *CanonicalizedQuadIterator {
	if !c.hasCanonicalQuad {
		panic("canonical quads was not enabled")
	}

	return &CanonicalizedQuadIterator{
		i: c.NewIterator(),
	}
}

func (c *Canonicalization) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	for _, nquad := range c.nquads {
		n, err := w.Write(nquad.encoded)
		totalWritten += int64(n)
		if err != nil {
			return totalWritten, err
		}
	}

	return totalWritten, nil
}

//

type CanonicalizedIterator struct {
	c   *Canonicalization
	idx int
}

func (it *CanonicalizedIterator) Close() error {
	return nil
}

func (it *CanonicalizedIterator) Err() error {
	return nil
}

func (it *CanonicalizedIterator) Next() bool {
	it.idx++

	return it.idx < len(it.c.nquads)
}

func (it *CanonicalizedIterator) OriginalQuadIndex() int64 {
	return it.c.nquads[it.idx].originalIndex
}

func (it *CanonicalizedIterator) EncodedQuad() []byte {
	return it.c.nquads[it.idx].encoded
}

//

type CanonicalizedQuadIterator struct {
	i *CanonicalizedIterator
}

var _ rdf.QuadIterator = &CanonicalizedQuadIterator{}

func (it *CanonicalizedQuadIterator) Close() error {
	return it.i.Close()
}

func (it *CanonicalizedQuadIterator) Err() error {
	return it.i.Err()
}

func (it *CanonicalizedQuadIterator) Next() bool {
	return it.i.Next()
}

func (it *CanonicalizedQuadIterator) Quad() rdf.Quad {
	return *it.i.c.nquads[it.i.idx].canonical
}

func (it *CanonicalizedQuadIterator) Statement() rdf.Statement {
	return it.Quad()
}

func (it *CanonicalizedQuadIterator) OriginalQuadIndex() int64 {
	return it.i.OriginalQuadIndex()
}

func (it *CanonicalizedQuadIterator) EncodedQuad() []byte {
	return it.i.EncodedQuad()
}
