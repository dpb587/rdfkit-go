package simplequery

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
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
	termsByVar          map[string]rdf.Term
	tripleBindingsByVar map[string]rdfio.Statement
}

func NewQueryResultBinding(termsByVar map[string]rdf.Term, tripleBindingsByVar map[string]rdfio.Statement) QueryResultBinding {
	return QueryResultBinding{
		termsByVar:          termsByVar,
		tripleBindingsByVar: tripleBindingsByVar,
	}
}

func (b QueryResultBinding) Clone() QueryResultBinding {
	termsByVar := map[string]rdf.Term{}
	tripleBindingsByVar := map[string]rdfio.Statement{}

	for k, v := range b.termsByVar {
		termsByVar[k] = v
	}

	for k, v := range b.tripleBindingsByVar {
		tripleBindingsByVar[k] = v
	}

	return QueryResultBinding{
		termsByVar:          termsByVar,
		tripleBindingsByVar: tripleBindingsByVar,
	}
}

func (b QueryResultBinding) Get(v string) rdf.Term {
	return b.termsByVar[v]
}

func (b QueryResultBinding) GetTripleBinding(v string) rdfio.Statement {
	return b.tripleBindingsByVar[v]
}
