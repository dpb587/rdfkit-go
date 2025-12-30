package quads

import "github.com/dpb587/rdfkit-go/rdf"

type TripleMatcher struct {
	Matcher rdf.TripleMatcher
}

var _ rdf.QuadMatcher = TripleMatcher{}

func (m TripleMatcher) MatchQuad(q rdf.Quad) bool {
	return m.Matcher.MatchTriple(q.Triple)
}
