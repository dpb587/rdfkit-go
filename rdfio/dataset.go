package rdfio

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

type DatasetReader interface {
	NewGraphIterator(ctx context.Context) (GraphIterator, error)
	NewGraphNameIterator(ctx context.Context) (GraphNameIterator, error)
	NewNodeIterator(ctx context.Context) (NodeIterator, error)
	NewStatementIterator(ctx context.Context, matchers ...StatementMatcher) (DatasetStatementIterator, error)
}

type Dataset interface {
	DatasetReader
	DatasetWriter

	GetGraph(ctx context.Context, graphName rdf.GraphNameValue) Graph

	GetStatement(ctx context.Context, triple rdf.Triple) (Statement, error)
	GetGraphStatement(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) (Statement, error)

	DeleteTriple(ctx context.Context, triple rdf.Triple) error
	DeleteGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error
}

type DatasetWriter interface {
	GraphWriter

	PutGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error
}
