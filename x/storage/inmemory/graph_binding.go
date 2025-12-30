package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/quads"
	"github.com/dpb587/rdfkit-go/rdf/triples"
)

type Graph struct {
	d     *Dataset
	t     rdf.GraphNameValue
	tNode *Node

	Baggage map[any]any

	assertedBySubject map[*Node]statementList
	// assertedByPredicate             map[*Node]statementList
	// assertedByObjectSubject         map[*Node]statementList
	// assertedByObjectLiteralDatatype map[rdf.IRI]statementList

	// annotationsByStatement map[*BoundStatement]BoundStatementList
}

func (gb *Graph) GetGraphName() rdf.GraphNameValue {
	return gb.t
}

func (gb *Graph) GetDataset() *Dataset {
	return gb.d
}

func (gb *Graph) newStatementIterator(matchers ...rdf.QuadMatcher) (*StatementIterator, error) {
	iter := &StatementIterator{
		index: -1,
	}

	if len(matchers) == 0 {
		for _, edges := range gb.assertedBySubject {
			iter.edges = append(iter.edges, edges...)
		}
	} else {
		var subjectMatchers []triples.SubjectMatcher
		var otherMatchers []rdf.QuadMatcher

		for _, matcher := range matchers {
			switch matcherType := matcher.(type) {
			case quads.TripleMatcher:
				switch submatcherType := matcherType.Matcher.(type) {
				case triples.SubjectMatcher:
					subjectMatchers = append(subjectMatchers, submatcherType)
				default:
					otherMatchers = append(otherMatchers, matcher)
				}
			default:
				otherMatchers = append(otherMatchers, matcher)
			}
		}

		if len(subjectMatchers) == 1 {
			for node := range gb.assertedBySubject {
				if !subjectMatchers[0].Matcher.MatchTerm(node.GetTerm()) {
					continue
				}

				for _, edge := range gb.assertedBySubject[node] {
					for _, matcher := range otherMatchers {
						if !matcher.MatchQuad(edge.GetQuad()) {
							goto NEXT_SUBJECT
						}
					}

					iter.edges = append(iter.edges, edge)

				NEXT_SUBJECT:
				}
			}
		} else {
			for _, edges := range gb.assertedBySubject {
				for _, edge := range edges {
					for _, matcher := range matchers {
						if !matcher.MatchQuad(edge.GetQuad()) {
							goto NEXT_EDGE
						}
					}

					iter.edges = append(iter.edges, edge)

				NEXT_EDGE:
				}
			}
		}
	}

	return iter, nil
}
