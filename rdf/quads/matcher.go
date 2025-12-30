package quads

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type GraphNameMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.QuadMatcher = GraphNameMatcher{}

func (m GraphNameMatcher) MatchQuad(t rdf.Quad) bool {
	return m.Matcher.MatchTerm(t.GraphName)
}

//

type SubjectMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.QuadMatcher = SubjectMatcher{}

func (m SubjectMatcher) MatchQuad(t rdf.Quad) bool {
	return m.Matcher.MatchTerm(t.Triple.Subject)
}

//

type PredicateMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.QuadMatcher = PredicateMatcher{}

func (m PredicateMatcher) MatchQuad(t rdf.Quad) bool {
	return m.Matcher.MatchTerm(t.Triple.Predicate)
}

//

type ObjectMatcher struct {
	Matcher rdf.TermMatcher
}

var _ rdf.QuadMatcher = ObjectMatcher{}

func (m ObjectMatcher) MatchQuad(t rdf.Quad) bool {
	return m.Matcher.MatchTerm(t.Triple.Object)
}

//

type TripleMatcher struct {
	Matcher rdf.TripleMatcher
}

var _ rdf.QuadMatcher = TripleMatcher{}

func (m TripleMatcher) MatchQuad(q rdf.Quad) bool {
	return m.Matcher.MatchTriple(q.Triple)
}
