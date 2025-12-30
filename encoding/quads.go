package encoding

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/quads"
)

type QuadsDecoder interface {
	rdf.QuadIterator
}

type QuadsEncoder interface {
	quads.DatasetWriter

	Close() error
	GetContentMetadata() ContentMetadata
}

type QuadsFactory interface {
	NewDecoder(r io.Reader) (QuadsDecoder, error)
	NewEncoder(w io.Writer) (QuadsEncoder, error)
	GetContentMetadata() ContentMetadata
}
