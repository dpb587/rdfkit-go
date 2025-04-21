package rdfioutil

import "github.com/dpb587/rdfkit-go/rdfio"

// TODO rename or remove

type GraphStatementMatcherIterator struct {
	parent  rdfio.StatementIterator
	matcher rdfio.StatementMatcher
}

var _ rdfio.StatementIterator = &GraphStatementMatcherIterator{}

func NewStatementMatcherIterator(parent rdfio.StatementIterator, matcher rdfio.StatementMatcher) *GraphStatementMatcherIterator {
	return &GraphStatementMatcherIterator{
		parent:  parent,
		matcher: matcher,
	}
}

func (i *GraphStatementMatcherIterator) Close() error {
	return i.parent.Close()
}

func (i *GraphStatementMatcherIterator) Err() error {
	return i.parent.Err()
}

func (i *GraphStatementMatcherIterator) Next() bool {
	for i.parent.Next() {
		if i.matcher == nil || i.matcher.MatchStatement(i.parent.GetStatement()) {
			return true
		}
	}

	return false
}

func (i *GraphStatementMatcherIterator) GetStatement() rdfio.Statement {
	return i.parent.GetStatement()
}
