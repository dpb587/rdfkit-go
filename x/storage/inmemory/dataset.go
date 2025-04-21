package inmemory

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

type GraphStatementImportHookFunc func(ctx context.Context, gb *Graph, tb *Statement, src rdfio.Statement) error
type GraphStatementHookFunc func(ctx context.Context, gb *Graph, tb *Statement)

type DatasetHooks struct {
	InitDataset   func(d *Dataset)
	InitGraph     func(tb *Graph)
	InitNode      func(tb *Node)
	InitStatement func(gb *Graph, tb *Statement)

	HandleImportStatement GraphStatementImportHookFunc

	StatementAdded   GraphStatementHookFunc
	StatementDeleted GraphStatementHookFunc
}

func WithDatasetHookManager(e DatasetHooks) DatasetOption {
	return datasetOptionFunc(func(d *Dataset) {
		d.hooks = e
	})
}

type Dataset struct {
	hooks DatasetHooks

	nodesByIRI          map[rdf.IRI]*Node
	nodesByBlankNodeRef map[rdf.BlankNode]*Node
	nodesByLiteral      map[[12]byte]*Node

	graphs map[rdf.GraphNameValue]*Graph

	Baggage map[any]any
}

var _ rdfio.Dataset = &Dataset{}

func NewDataset(opts ...DatasetOption) *Dataset {
	d := &Dataset{
		nodesByIRI:          map[rdf.IRI]*Node{},
		nodesByBlankNodeRef: map[rdf.BlankNode]*Node{},
		nodesByLiteral:      map[[12]byte]*Node{},
		graphs:              map[rdf.GraphNameValue]*Graph{},
	}

	for _, opt := range opts {
		opt.apply(d)
	}

	if d.hooks.InitDataset != nil {
		d.hooks.InitDataset(d)
	}

	d.graphs[rdf.DefaultGraph] = d.createGraph(rdf.DefaultGraph)

	return d
}

func (d *Dataset) Close() error {
	return nil
}

func (d *Dataset) createGraph(graphName rdf.GraphNameValue) *Graph {
	var tNode *Node

	if graphName != rdf.DefaultGraph {
		tNode, _ = d.bindNode(graphName.(rdf.Term), true)
	}

	res := &Graph{
		d:                 d,
		t:                 graphName,
		tNode:             tNode,
		assertedBySubject: map[*Node]statementList{},
		// assertedByPredicate:             map[*Node]statementList{},
		// assertedByObjectSubject:         map[*Node]statementList{},
		// assertedByObjectLiteralDatatype: map[rdf.IRI]statementList{},
	}

	if d.hooks.InitGraph != nil {
		d.hooks.InitGraph(res)
	}

	d.graphs[graphName] = res

	return res
}

func (d *Dataset) NewGraphNameIterator(ctx context.Context) (rdfio.GraphNameIterator, error) {
	var all rdf.GraphNameValueList

	for graphName, graph := range d.graphs {
		if len(graph.assertedBySubject) == 0 {
			continue
		}

		all = append(all, graphName)
	}

	return rdfioutil.NewStaticGraphNameIterator(all), nil
}

func (d *Dataset) NewGraphIterator(ctx context.Context) (rdfio.GraphIterator, error) {
	var all []rdfio.Graph

	for _, graph := range d.graphs {
		if len(graph.assertedBySubject) == 0 {
			continue
		}

		all = append(all, graph)
	}

	return rdfioutil.NewStaticGraphIterator(all), nil
}

func (d *Dataset) PutTriple(ctx context.Context, triple rdf.Triple) error {
	return d.graphPutTriple(ctx, d.graphs[rdf.DefaultGraph], triple, nil)
}

func (d *Dataset) PutGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	graph, ok := d.graphs[graphName]
	if !ok {
		graph = d.createGraph(graphName)
	}

	return d.graphPutTriple(ctx, graph, triple, nil)
}

// TODO consolidate Import vs WriteStatement

func (d *Dataset) ImportStatement(ctx context.Context, s rdfio.Statement) error {
	graphName := s.GetGraphName()

	graph, ok := d.graphs[graphName]
	if !ok {
		graph = d.createGraph(graphName)
	}

	return d.graphPutTriple(ctx, graph, s.GetTriple(), func(ctx context.Context, gb *Graph, tb *Statement) {
		if d.hooks.HandleImportStatement != nil {
			err := d.hooks.HandleImportStatement(ctx, gb, tb, s)
			if err != nil {
				panic(fmt.Sprintf("import statement: %v", err)) // TODO propagate
			}
		}
	})
}

func (d *Dataset) graphPutTriple(ctx context.Context, graph *Graph, triple rdf.Triple, hook GraphStatementHookFunc) error {
	statement, exists, err := d.bindStatement(graph, triple)
	if err != nil {
		return fmt.Errorf("bind triple: %v", err)
	} else if exists {
		if hook != nil {
			hook(ctx, graph, statement)
		} else if d.hooks.StatementAdded != nil {
			d.hooks.StatementAdded(ctx, graph, statement)
		}

		return nil
	}

	if d.hooks.InitStatement != nil {
		d.hooks.InitStatement(graph, statement)
	}

	graph.assertedBySubject[statement.s] = append(graph.assertedBySubject[statement.s], statement)
	// graph.assertedByPredicate[statement.p] = append(graph.assertedByPredicate[statement.p], statement)

	// switch t := statement.o.t.(type) {
	// case rdf.BlankNode, rdf.IRI:
	// 	graph.assertedByObjectSubject[statement.o] = append(graph.assertedByObjectSubject[statement.o], statement)
	// case rdf.Literal:
	// 	graph.assertedByObjectLiteralDatatype[t.Datatype] = append(graph.assertedByObjectLiteralDatatype[t.Datatype], statement)
	// }

	if hook != nil {
		hook(ctx, graph, statement)
	} else if d.hooks.StatementAdded != nil {
		d.hooks.StatementAdded(ctx, graph, statement)
	}

	return nil
}

func (d *Dataset) DeleteTriple(ctx context.Context, triple rdf.Triple) error {
	return d.graphDeleteTriple(ctx, d.graphs[rdf.DefaultGraph], triple)
}

func (d *Dataset) DeleteGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	graph, ok := d.graphs[graphName]
	if !ok {
		return nil
	}

	return d.graphDeleteTriple(ctx, graph, triple)
}

func (d *Dataset) graphDeleteTriple(ctx context.Context, graph *Graph, triple rdf.Triple) error {
	statement, exists, err := d.bindStatement(graph, triple)
	if err != nil {
		return fmt.Errorf("bind triple: %v", err)
	} else if !exists {
		return nil
	}

	if d.hooks.StatementDeleted != nil {
		d.hooks.StatementDeleted(ctx, graph, statement)
	}

	graph.assertedBySubject[statement.s] = graph.assertedBySubject[statement.s].Exclude(statement)
	// graph.assertedByPredicate[statement.p] = graph.assertedByPredicate[statement.p].Exclude(statement)

	// switch t := statement.o.t.(type) {
	// case rdf.BlankNode, rdf.IRI:
	// 	graph.assertedByObjectSubject[statement.o] = graph.assertedByObjectSubject[statement.o].Exclude(statement)
	// case rdf.Literal:
	// 	graph.assertedByObjectLiteralDatatype[t.Datatype] = graph.assertedByObjectLiteralDatatype[t.Datatype].Exclude(statement)
	// }

	return nil
}

func (d *Dataset) GetNode(ctx context.Context, s rdf.SubjectValue) (rdfio.Node, error) {
	return d.getNode(ctx, d.graphs[rdf.DefaultGraph], s)
}

func (d *Dataset) getNode(_ context.Context, graph *Graph, s rdf.SubjectValue) (rdfio.Node, error) {
	boundSubject, known := d.bindNode(s, false)
	if !known {
		return nil, rdfio.ErrNodeNotBound
	} else if len(graph.assertedBySubject[boundSubject]) == 0 {
		return nil, rdfio.ErrNodeNotBound
	}

	return boundSubject, nil
}

func (d *Dataset) NewNodeIterator(ctx context.Context, matchers ...rdfio.StatementMatcher) (rdfio.DatasetNodeIterator, error) {
	panic("TODO")
}

func (d *Dataset) GetStatement(ctx context.Context, triple rdf.Triple) (rdfio.Statement, error) {
	return d.getStatement(ctx, d.graphs[rdf.DefaultGraph], triple)
}

func (d *Dataset) GetGraph(ctx context.Context, graphName rdf.GraphNameValue) rdfio.Graph {
	graph, ok := d.graphs[graphName]
	if !ok {
		graph = d.createGraph(graphName)
	}

	return graph
}

func (d *Dataset) GetGraphStatement(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) (rdfio.Statement, error) {
	graph, ok := d.graphs[graphName]
	if !ok {
		return nil, rdfio.ErrStatementNotBound
	}

	return d.getStatement(ctx, graph, triple)
}

func (d *Dataset) getStatement(_ context.Context, graph *Graph, triple rdf.Triple) (rdfio.Statement, error) {
	statement, exists, err := d.bindStatement(graph, triple)
	if err != nil {
		return nil, fmt.Errorf("bind triple: %v", err)
	} else if !exists {
		return nil, rdfio.ErrStatementNotBound
	}

	return statement, nil
}

func (d *Dataset) NewStatementIterator(ctx context.Context, matchers ...rdfio.StatementMatcher) (rdfio.DatasetStatementIterator, error) {
	// TODO this should be all graphs
	return d.graphs[rdf.DefaultGraph].newStatementIterator(matchers...)
}

func (d *Dataset) bindNode(t rdf.Term, write bool) (*Node, bool) {
	switch t := t.(type) {
	case rdf.IRI:
		nb, ok := d.nodesByIRI[t]
		if ok {
			return nb, true
		} else if !write {
			return nil, false
		}

		nb = &Node{
			d: d,
			t: t,
		}

		if d.hooks.InitNode != nil {
			d.hooks.InitNode(nb)
		}

		d.nodesByIRI[t] = nb

		return nb, true
	case rdf.BlankNode:
		nb, ok := d.nodesByBlankNodeRef[t]
		if ok {
			return nb, true
		} else if !write {
			return nil, false
		}

		nb = &Node{
			d: d,
			t: t,
		}

		if d.hooks.InitNode != nil {
			d.hooks.InitNode(nb)
		}

		d.nodesByBlankNodeRef[t] = nb

		return nb, true
	case rdf.Literal:
		h := sha256.New()

		h.Write([]byte(t.Datatype + "\n"))

		for k, v := range t.Tags {
			fmt.Fprintf(h, "%v=%q\n", k, v)
		}

		h.Write([]byte(t.LexicalForm))

		key := [12]byte(h.Sum(nil)[0:12])

		nb, ok := d.nodesByLiteral[key]
		if ok {
			return nb, true
		} else if !write {
			return nil, false
		}

		nb = &Node{
			d: d,
			t: t,
		}

		if d.hooks.InitNode != nil {
			d.hooks.InitNode(nb)
		}

		d.nodesByLiteral[key] = nb

		return nb, true
	}

	panic(fmt.Sprintf("unsupported node type: %T", t))
}

func (d *Dataset) bindStatement(boundGraph *Graph, triple rdf.Triple) (*Statement, bool, error) {
	boundSubject, _ := d.bindNode(triple.Subject, true)     // TODO presumptive write
	boundPredicate, _ := d.bindNode(triple.Predicate, true) // TODO presumptive write
	boundObject, _ := d.bindNode(triple.Object, true)       // TODO presumptive write

	for _, known := range boundGraph.assertedBySubject[boundSubject] {
		if known.p != boundPredicate {
			continue
		} else if known.o != boundObject {
			continue
		}

		return known, true, nil
	}

	statement := &Statement{
		g: boundGraph,
		s: boundSubject,
		p: boundPredicate,
		o: boundObject,
	}

	return statement, false, nil
}
