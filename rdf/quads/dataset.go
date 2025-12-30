package quads

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

type DatasetWriter interface {
	AddQuad(ctx context.Context, quad rdf.Quad) error
}

type Dataset interface {
	DatasetWriter

	NewQuadIterator(ctx context.Context, matchers ...rdf.QuadMatcher) (rdf.QuadIterator, error)

	HasQuad(ctx context.Context, quad rdf.Quad) (bool, error)
	DeleteQuad(ctx context.Context, quad rdf.Quad) error
}
