package triples

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

type GraphWriter interface {
	AddTriple(ctx context.Context, triple rdf.Triple) error
}

type Graph interface {
	GraphWriter

	NewTripleIterator(ctx context.Context, matchers ...rdf.TripleMatcher) (rdf.TripleIterator, error)

	HasTriple(ctx context.Context, triple rdf.Triple) (bool, error)
	DeleteTriple(ctx context.Context, triple rdf.Triple) error
}
