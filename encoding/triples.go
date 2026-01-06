package encoding

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/triples"
)

type TriplesDecoder interface {
	Decoder
	rdf.TripleIterator
}

type TriplesEncoder interface {
	Encoder
	triples.GraphWriter
}

type TriplesFactory interface {
	NewDecoder(r io.Reader) (TriplesDecoder, error)
	NewEncoder(w io.Writer) (TriplesEncoder, error)
	GetContentMetadata() ContentMetadata
}
