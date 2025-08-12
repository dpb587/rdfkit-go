package inmemory

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
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

var _ rdfio.Graph = &Graph{}

func (gb *Graph) Close() error {
	// TODO rdfio.DatasetGraphRepository without Close?
	panic("not implemented for graph of dataset")
}

func (gb *Graph) GetGraphName() rdf.GraphNameValue {
	return gb.t
}

func (gb *Graph) GetGraphNameNode() rdfio.Node {
	return gb.tNode
}

func (gb *Graph) NewStatementIterator(ctx context.Context, matchers ...rdfio.StatementMatcher) (rdfio.GraphStatementIterator, error) {
	iter, err := gb.newStatementIterator(matchers...)

	return iter, err
}

func (gb *Graph) NewNodeIterator(ctx context.Context) (rdfio.NodeIterator, error) {
	v, err := gb.newNodeIterator()

	return v, err
}

func (gb *Graph) GetDataset() rdfio.Dataset {
	return gb.d
}

func (gb *Graph) GetStatement(ctx context.Context, triple rdf.Triple) (rdfio.Statement, error) {
	return gb.d.getStatement(ctx, gb, triple)
}

func (gb *Graph) GetNode(ctx context.Context, s rdf.SubjectValue) (rdfio.Node, error) {
	return gb.d.GetNode(ctx, s)
}

func (gb *Graph) DeleteTriple(ctx context.Context, triple rdf.Triple) error {
	return gb.d.graphDeleteTriple(ctx, gb, triple)
}

func (gb *Graph) PutTriple(ctx context.Context, triple rdf.Triple) error {
	return gb.d.graphPutTriple(ctx, gb, triple, nil)
}

func (gb *Graph) newNodeIterator() (rdfio.NodeIterator, error) {
	nodes := make(nodeList, 0, len(gb.assertedBySubject))

	for node := range gb.assertedBySubject {
		nodes = append(nodes, node)
	}

	return &nodeIterator{
		index: -1,
		edges: nodes,
	}, nil
}

func (gb *Graph) newStatementIterator(matchers ...rdfio.StatementMatcher) (*statementIterator, error) {
	iter := &statementIterator{
		index: -1,
	}

	if len(matchers) == 0 {
		for _, edges := range gb.assertedBySubject {
			iter.edges = append(iter.edges, edges...)
		}
	} else {
		var subjectMatchers []rdfio.SubjectStatementMatcher
		var otherMatchers []rdfio.StatementMatcher

		for _, matcher := range matchers {
			switch matcherType := matcher.(type) {
			case rdfio.SubjectStatementMatcher:
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
						if !matcher.MatchStatement(edge) {
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
						if !matcher.MatchStatement(edge) {
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
