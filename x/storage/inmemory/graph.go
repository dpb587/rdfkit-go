package inmemory

import (
	"context"

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

var _ triples.Graph = &Graph{}

func (gb *Graph) GetGraphName() rdf.GraphNameValue {
	return gb.t
}

func (gb *Graph) GetDataset() *Dataset {
	return gb.d
}

func (gb *Graph) AddTriple(ctx context.Context, t rdf.Triple) error {
	return gb.d.addQuad(ctx, rdf.Quad{
		Triple:    t,
		GraphName: gb.t,
	})
}

func (gb *Graph) DeleteTriple(ctx context.Context, t rdf.Triple) error {
	return gb.d.DeleteQuad(ctx, rdf.Quad{
		Triple:    t,
		GraphName: gb.t,
	})
}

func (gb *Graph) HasTriple(ctx context.Context, t rdf.Triple) (bool, error) {
	return gb.d.HasQuad(ctx, rdf.Quad{
		Triple:    t,
		GraphName: gb.t,
	})
}

func (gb *Graph) GetTripleStatement(ctx context.Context, triple rdf.Triple) (*Statement, error) {
	return gb.d.GetQuadStatement(ctx, rdf.Quad{
		Triple:    triple,
		GraphName: gb.t,
	})
}

func (gb *Graph) NewTripleIterator(ctx context.Context, matchers ...rdf.TripleMatcher) (rdf.TripleIterator, error) {
	iter := &TripleIterator{
		index: -1,
	}

	if len(matchers) == 0 {
		for _, edges := range gb.assertedBySubject {
			iter.edges = append(iter.edges, edges...)
		}
	} else {
		var subjectMatchers []triples.SubjectMatcher
		var otherMatchers []rdf.TripleMatcher

		for _, matcher := range matchers {
			switch matcherType := matcher.(type) {
			case triples.SubjectMatcher:
				subjectMatchers = append(subjectMatchers, matcherType)
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
						if !matcher.MatchTriple(edge.GetQuad().Triple) {
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
						if !matcher.MatchTriple(edge.GetQuad().Triple) {
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

func (gb *Graph) newQuadIterator(matchers ...rdf.QuadMatcher) (*QuadIterator, error) {
	iter := &QuadIterator{
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

func (gb *Graph) NewSubjectIterator(ctx context.Context, matchers ...rdf.TermMatcher) (rdf.TermIterator, error) {
	iter := &termIterator{
		index: -1,
	}

	for node := range gb.assertedBySubject {
		for _, matcher := range matchers {
			if !matcher.MatchTerm(node.GetTerm()) {
				goto NEXT_SUBJECT
			}
		}

		iter.edges = append(iter.edges, node)

	NEXT_SUBJECT:
	}

	return iter, nil
}
