package encodingutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type QuadAsTripleEncoder struct {
	encoding.TriplesEncoder
}

var _ encoding.QuadsEncoder = &QuadAsTripleEncoder{}
var _ encoding.TriplesEncoder = &QuadAsTripleEncoder{}

func (e QuadAsTripleEncoder) AddQuad(ctx context.Context, quad rdf.Quad) error {
	return e.TriplesEncoder.AddTriple(ctx, quad.Triple)
}
