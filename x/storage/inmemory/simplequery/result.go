package simplequery

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type QueryResult interface {
	Close() error
	Err() error
	Next() bool
	GetBinding() QueryResultBinding
}

type queryResultBindingIterator struct {
	bindings []QueryResultBinding
	index    int
}

// var _ QueryResult = &queryResultBindingIterator{}

func NewQueryResult(bindings []QueryResultBinding) QueryResult {
	return &queryResultBindingIterator{
		bindings: bindings,
		index:    -1,
	}
}

func (i *queryResultBindingIterator) Close() error {
	return nil
}

func (i *queryResultBindingIterator) Err() error {
	return nil
}

func (i *queryResultBindingIterator) Next() bool {
	if i.index >= len(i.bindings)-1 {
		return false
	}

	i.index++

	return true
}

func (i *queryResultBindingIterator) GetBinding() QueryResultBinding {
	return i.bindings[i.index]
}

type QueryResultBinding struct {
	termsByVar map[string]rdf.Term
}

func NewQueryResultBinding(termsByVar map[string]rdf.Term) QueryResultBinding {
	return QueryResultBinding{
		termsByVar: termsByVar,
	}
}

func (b QueryResultBinding) Clone() QueryResultBinding {
	termsByVar := map[string]rdf.Term{}

	for k, v := range b.termsByVar {
		termsByVar[k] = v
	}

	return QueryResultBinding{
		termsByVar: termsByVar,
	}
}

func (b QueryResultBinding) Get(v string) rdf.Term {
	return b.termsByVar[v]
}
