package encodingutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type quadIteratorIterator struct {
	iters        []rdf.QuadIterator
	textOffsets0 encoding.StatementTextOffsetsProvider
	err          error
}

var _ encoding.QuadsDecoder = &quadIteratorIterator{}
var _ encoding.StatementTextOffsetsProvider = &quadIteratorIterator{}

func NewQuadIteratorIterator(iters ...rdf.QuadIterator) encoding.QuadsDecoder {
	i := &quadIteratorIterator{
		iters: iters,
	}
	i.iter0()

	return i
}

func (i *quadIteratorIterator) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return "public.composite-content"
}

func (i *quadIteratorIterator) Close() error {
	for _, iter := range i.iters {
		iter.Close()
	}

	return nil
}

func (i *quadIteratorIterator) Err() error {
	return i.err
}

func (i *quadIteratorIterator) Next() bool {
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
		i.iter0()
	}
}

func (i *quadIteratorIterator) Quad() rdf.Quad {
	return i.iters[0].Quad()
}

func (i *quadIteratorIterator) Statement() rdf.Statement {
	return i.Quad()
}

func (i *quadIteratorIterator) StatementTextOffsets() encoding.StatementTextOffsets {
	if i.textOffsets0 == nil {
		return nil
	}

	return i.textOffsets0.StatementTextOffsets()
}

func (i *quadIteratorIterator) iter0() {
	if len(i.iters) == 0 {
		i.textOffsets0 = nil

		return
	}

	i.textOffsets0, _ = i.iters[0].(encoding.StatementTextOffsetsProvider)
}
