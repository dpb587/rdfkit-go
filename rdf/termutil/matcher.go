package termutil

import "github.com/dpb587/rdfkit-go/rdf"

type Matcher interface {
	MatchTerm(t rdf.Term) bool
}

type MatcherFunc func(t rdf.Term) bool

var _ Matcher = MatcherFunc(nil)

func (f MatcherFunc) MatchTerm(t rdf.Term) bool {
	return f(t)
}
