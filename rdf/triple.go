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

// TODO needs more research; subject+object? wait for next version of the spec?
// in the meantime, standalone struct is meaningful elsewhere and unlikely to change with new spec

// var _ Term = Triple{}
// var _ ObjectValue = Triple{}

// func (Triple) isTermBuiltin()        {}
// func (Triple) isObjectValueBuiltin() {}

// func (Triple) TermKind() TermKind {
// 	return TermKindTriple
// }

// func (t Triple) TermEquals(d Term) bool {
// 	dTriple, ok := d.(Triple)
// 	if !ok {
// 		return false
// 	} else if !t.Subject.TermEquals(dTriple.Subject) {
// 		return false
// 	} else if !t.Predicate.TermEquals(dTriple.Predicate) {
// 		return false
// 	}

// 	return t.Object.TermEquals(dTriple.Object)
// }

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
