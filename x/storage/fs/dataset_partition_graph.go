package fs

// import (
// 	"context"

// 	"github.com/dpb587/rdfkit-go/rdf"
// 	"github.com/dpb587/rdfkit-go/rdfio"
// )

// type datasetPartitionGraph struct {
// 	dp        *datasetPartition
// 	graphName rdf.GraphNameValue
// }

// var _ rdfio.Graph = &datasetPartitionGraph{}

// func (dpg *datasetPartitionGraph) NewNodeIterator(ctx context.Context) (rdfio.NodeIterator, error) {
// 	err := dpg.dp.requireLoad(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return dpg.dp.dataset.GetGraph(ctx, dpg.graphName).NewNodeIterator(ctx)
// }

// func (dpg *datasetPartitionGraph) NewStatementIterator(ctx context.Context, matchers ...rdfio.StatementMatcher) (rdfio.GraphStatementIterator, error) {
// 	err := dpg.dp.requireLoad(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return dpg.dp.dataset.GetGraph(ctx, dpg.graphName).NewStatementIterator(ctx, matchers...)
// }

// func (dpg *datasetPartitionGraph) GetGraphName() rdf.GraphNameValue {
// 	return dpg.graphName
// }

// func (dpg *datasetPartitionGraph) GetNode(ctx context.Context, t rdf.SubjectValue) (rdfio.Node, error) {
// 	err := dpg.dp.requireLoad(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return dpg.dp.dataset.GetGraph(ctx, dpg.graphName).GetNode(ctx, t)
// }

// func (dpg *datasetPartitionGraph) GetStatement(ctx context.Context, triple rdf.Triple) (rdfio.Statement, error) {
// 	err := dpg.dp.requireLoad(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return dpg.dp.dataset.GetGraph(ctx, dpg.graphName).GetStatement(ctx, triple)
// }

// func (dpg *datasetPartitionGraph) DeleteTriple(ctx context.Context, triple rdf.Triple) error {
// 	err := dpg.dp.requireLoad(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	dpg.dp.isModified = true

// 	return dpg.dp.dataset.GetGraph(ctx, dpg.graphName).DeleteTriple(ctx, triple)
// }

// func (dpg *datasetPartitionGraph) Close() error {
// 	return dpg.dp.Close()
// }

// func (dpg *datasetPartitionGraph) PutTriple(ctx context.Context, triple rdf.Triple) error {
// 	err := dpg.dp.requireLoad(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	dpg.dp.isModified = true

// 	return dpg.dp.dataset.GetGraph(ctx, dpg.graphName).PutTriple(ctx, triple)
// }
