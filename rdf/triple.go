package rdf

type Triple struct {
	Subject   SubjectValue
	Predicate PredicateValue
	Object    ObjectValue
}

var _ Statement = Triple{}

func (Triple) isStatement() {}

func (t Triple) StatementType() StatementType {
	return TripleStatementType
}

func (t Triple) AsQuad(g GraphNameValue) Quad {
	return Quad{
		Triple:    t,
		GraphName: g,
	}
}

//

type TripleMatcher interface {
	MatchTriple(t Triple) bool
}

//

type TripleMatcherFunc func(t Triple) bool

func (f TripleMatcherFunc) MatchTriple(t Triple) bool {
	return f(t)
}

//

type TripleList []Triple

func (tl TripleList) AsQuads(g GraphNameValue) QuadList {
	quads := make(QuadList, 0, len(tl))

	for _, t := range tl {
		quads = append(quads, t.AsQuad(g))
	}

	return quads
}

//

type TripleIterator interface {
	StatementIterator

	Triple() Triple
}
