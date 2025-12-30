package inmemory

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/quads"
)

var ErrNoStatement = fmt.Errorf("no statement")

type StatementHook func(ctx context.Context, tb *Statement)

type GraphStatementHookFunc func(ctx context.Context, gb *Graph, tb *Statement)

type DatasetHooks struct {
	InitDataset   func(d *Dataset)
	InitGraph     func(tb *Graph)
	InitNode      func(tb *Node)
	InitStatement func(gb *Graph, tb *Statement)

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

var _ quads.Dataset = &Dataset{}

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

	d.graphs[nil] = d.createGraph(nil)

	return d
}

func (d *Dataset) Close() error {
	return nil
}

func (d *Dataset) createGraph(graphName rdf.GraphNameValue) *Graph {
	var tNode *Node

	if graphName != nil {
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

func (d *Dataset) AddQuad(ctx context.Context, quad rdf.Quad) error {
	return d.addQuad(ctx, quad, nil)
}

func (d *Dataset) AddQuadStatement(ctx context.Context, quad rdf.Quad, f StatementHook) error {
	return d.addQuad(ctx, quad, f)
}

func (d *Dataset) addQuad(ctx context.Context, quad rdf.Quad, f StatementHook) error {
	graph, ok := d.graphs[quad.GraphName]
	if !ok {
		graph = d.createGraph(quad.GraphName)
	}

	statement, exists, err := d.bindStatement(graph, quad)
	if err != nil {
		return fmt.Errorf("bind quad: %v", err)
	} else if exists {
		if f != nil {
			f(ctx, statement)
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

	if d.hooks.StatementAdded != nil {
		d.hooks.StatementAdded(ctx, graph, statement)
	}

	if f != nil {
		f(ctx, statement)
	}

	return nil
}

func (d *Dataset) DeleteQuad(ctx context.Context, quad rdf.Quad) error {
	graph, ok := d.graphs[quad.GraphName]
	if !ok {
		graph = d.createGraph(quad.GraphName)
	}

	statement, exists, err := d.bindStatement(graph, quad)
	if err != nil {
		return fmt.Errorf("bind quad: %v", err)
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

func (d *Dataset) HasQuad(ctx context.Context, quad rdf.Quad) (bool, error) {
	_, err := d.GetStatement(ctx, quad)
	if err != nil {
		if err == ErrNoStatement {
			return false, nil
		}

		return false, fmt.Errorf("get quad: %v", err)
	}

	return true, nil
}

func (d *Dataset) GetStatement(ctx context.Context, quad rdf.Quad) (*Statement, error) {
	graph, ok := d.graphs[quad.GraphName]
	if !ok {
		return nil, ErrNoStatement
	}

	statement, exists, err := d.bindStatement(graph, quad)
	if err != nil {
		return nil, fmt.Errorf("bind quad: %v", err)
	} else if !exists {
		return nil, ErrNoStatement
	}

	return statement, nil
}

func (d *Dataset) NewQuadIterator(ctx context.Context, matchers ...rdf.QuadMatcher) (rdf.QuadIterator, error) {
	var all statementList

	for _, g := range d.graphs {
		iter, err := g.newStatementIterator(matchers...)
		if err != nil {
			return nil, fmt.Errorf("create graph statement iterator: %v", err)
		}

		all = append(all, iter.edges...)
	}

	return &StatementIterator{
		edges: all,
		index: -1,
	}, nil
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

		if t.Tag != nil {
			switch tag := t.Tag.(type) {
			case rdf.LanguageLiteralTag:
				fmt.Fprintf(h, "lang=%q\n", tag.Language)
			case rdf.DirectionalLanguageLiteralTag:
				fmt.Fprintf(h, "lang=%q; dir=%q\n", tag.Language, tag.BaseDirection)
			}
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

func (d *Dataset) bindStatement(boundGraph *Graph, quad rdf.Quad) (*Statement, bool, error) {
	boundSubject, _ := d.bindNode(quad.Triple.Subject, true)     // TODO presumptive write
	boundPredicate, _ := d.bindNode(quad.Triple.Predicate, true) // TODO presumptive write
	boundObject, _ := d.bindNode(quad.Triple.Object, true)       // TODO presumptive write

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
