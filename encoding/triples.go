package encoding

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/triples"
)

type TriplesDecoder interface {
	rdf.TripleIterator
}

type TriplesEncoder interface {
	triples.GraphWriter

	Close() error
	GetContentMetadata() ContentMetadata
}

type TriplesFactory interface {
	NewDecoder(r io.Reader) (TriplesDecoder, error)
	NewEncoder(w io.Writer) (TriplesEncoder, error)
	GetContentMetadata() ContentMetadata
}
