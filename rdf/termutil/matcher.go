package termutil

import "github.com/dpb587/rdfkit-go/rdf"

type Matcher interface {
	MatchTerm(t rdf.Term) bool
}

//

type MatcherFunc func(t rdf.Term) bool

var _ Matcher = MatcherFunc(nil)

func (f MatcherFunc) MatchTerm(t rdf.Term) bool {
	return f(t)
}

//

type LogicalOrMatcher []Matcher

var _ Matcher = LogicalOrMatcher(nil)

func (m LogicalOrMatcher) MatchTerm(t rdf.Term) bool {
	for _, matcher := range m {
		if matcher.MatchTerm(t) {
			return true
		}
	}

	return false
}

//

type LogicalAndMatcher []Matcher

var _ Matcher = LogicalAndMatcher(nil)

func (m LogicalAndMatcher) MatchTerm(t rdf.Term) bool {
	for _, matcher := range m {
		if !matcher.MatchTerm(t) {
			return false
		}
	}

	return true
}

//

type LogicalNotMatcher struct {
	Matcher Matcher
}

var _ Matcher = LogicalNotMatcher{}

func (m LogicalNotMatcher) MatchTerm(t rdf.Term) bool {
	return !m.Matcher.MatchTerm(t)
}
