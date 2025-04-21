package rdfio

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

type GraphReader interface {
	NewNodeIterator(ctx context.Context, matchers ...StatementMatcher) (GraphNodeIterator, error)
	NewStatementIterator(ctx context.Context, matchers ...StatementMatcher) (GraphStatementIterator, error)
}

type Graph interface {
	GraphReader
	GraphWriter

	GetGraphName() rdf.GraphNameValue

	GetNode(ctx context.Context, t rdf.SubjectValue) (Node, error)
	GetStatement(ctx context.Context, triple rdf.Triple) (Statement, error)

	DeleteTriple(ctx context.Context, triple rdf.Triple) error
}

type GraphWriter interface {
	Close() error

	PutTriple(ctx context.Context, triple rdf.Triple) error
}
