package rdfio

type GraphIterator interface {
	Close() error
	Err() error
	Next() bool
	GetGraph() Graph
}
