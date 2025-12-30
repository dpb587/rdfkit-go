package triples

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type SubjectMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.TripleMatcher = SubjectMatcher{}

func (m SubjectMatcher) MatchTriple(t rdf.Triple) bool {
	return m.Matcher.MatchTerm(t.Subject)
}

//

type PredicateMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.TripleMatcher = PredicateMatcher{}

func (m PredicateMatcher) MatchTriple(t rdf.Triple) bool {
	return m.Matcher.MatchTerm(t.Predicate)
}

//

type ObjectMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.TripleMatcher = ObjectMatcher{}

func (m ObjectMatcher) MatchTriple(t rdf.Triple) bool {
	return m.Matcher.MatchTerm(t.Object)
}
