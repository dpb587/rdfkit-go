package rdf

type Quad struct {
	Triple    Triple
	GraphName GraphNameValue
}

var _ Statement = Quad{}

func (Quad) isStatement() {}

func (q Quad) StatementType() StatementType {
	return QuadStatementType
}

//

type QuadMatcher interface {
	MatchQuad(t Quad) bool
}

//

type QuadList []Quad

func (ql QuadList) AsTriples() TripleList {
	triples := make(TripleList, 0, len(ql))

	for _, q := range ql {
		triples = append(triples, q.Triple)
	}

	return triples
}

//

type QuadIterator interface {
	StatementIterator

	Quad() Quad
}
