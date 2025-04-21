package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

// type StatementKind uint8

// const (
// 	StatementKindUnasserted StatementKind = iota
// 	StatementKindAnnotation
// 	StatementKindAsserted
// )

type Statement struct {
	g *Graph
	s *Node
	p *Node // not traditionally considered a node in RDF graphs; avoid exposing for now?
	o *Node

	Baggage map[any]any
}

var _ rdfio.Statement = &Statement{}

func (tb *Statement) GetGraphName() rdf.GraphNameValue {
	return tb.g.t.(rdf.GraphNameValue)
}

func (tb *Statement) GetTriple() rdf.Triple {
	return rdf.Triple{
		Subject:   tb.s.t.(rdf.SubjectValue),
		Predicate: tb.p.t.(rdf.PredicateValue),
		Object:    tb.o.t.(rdf.ObjectValue),
	}
}

func (tb *Statement) GetDataset() *Dataset {
	return tb.g.d
}

func (tb *Statement) GetGraph() *Graph {
	return tb.g
}

func (tb *Statement) GetSubjectNode() *Node {
	return tb.s
}

func (tb *Statement) GetObjectNode() *Node {
	return tb.o
}

//

type statementList []*Statement

func (p statementList) Exclude(s *Statement) statementList {
	var next statementList

	for _, statement := range p {
		if statement == s {
			continue
		}

		next = append(next, statement)
	}

	return next
}
