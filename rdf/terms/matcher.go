package terms

import "github.com/dpb587/rdfkit-go/rdf"

type LogicalOrMatcher []rdf.TermMatcher

var _ rdf.TermMatcher = LogicalOrMatcher(nil)

func (m LogicalOrMatcher) MatchTerm(t rdf.Term) bool {
	for _, matcher := range m {
		if matcher.MatchTerm(t) {
			return true
		}
	}

	return false
}

//

type LogicalAndMatcher []rdf.TermMatcher

var _ rdf.TermMatcher = LogicalAndMatcher(nil)

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
	Matcher rdf.TermMatcher
}

var _ rdf.TermMatcher = LogicalNotMatcher{}

func (m LogicalNotMatcher) MatchTerm(t rdf.Term) bool {
	return !m.Matcher.MatchTerm(t)
}
