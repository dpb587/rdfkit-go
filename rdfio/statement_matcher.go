package rdfio

import (
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type StatementMatcher interface {
	MatchStatement(e Statement) bool
}

//

type StatementMatcherFunc func(e Statement) bool

var _ StatementMatcher = StatementMatcherFunc(nil)

func (f StatementMatcherFunc) MatchStatement(e Statement) bool {
	return f(e)
}

//

type SubjectStatementMatcher struct {
	Matcher termutil.Matcher
}

func (m SubjectStatementMatcher) MatchStatement(e Statement) bool {
	return m.Matcher.MatchTerm(e.GetTriple().Subject)
}

//

type PredicateStatementMatcher struct {
	Matcher termutil.Matcher
}

func (m PredicateStatementMatcher) MatchStatement(e Statement) bool {
	return m.Matcher.MatchTerm(e.GetTriple().Predicate)
}

//

type ObjectStatementMatcher struct {
	Matcher termutil.Matcher
}

func (m ObjectStatementMatcher) MatchStatement(e Statement) bool {
	return m.Matcher.MatchTerm(e.GetTriple().Object)
}
