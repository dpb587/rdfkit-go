package fs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory"
)

type datasetPartition struct {
	dpm      *DatasetPartitionManager
	fileBase string
	dataset  *inmemory.Dataset

	isLoaded   bool
	isModified bool
}

var _ rdfio.Dataset = &datasetPartition{}

func (dp *datasetPartition) getFilePath() string {
	return filepath.Join(dp.dpm.dir, dp.fileBase+dp.dpm.documentEncoding.GetDatasetEncoderContentMetadata().FileExt)
}

func (dp *datasetPartition) Close() error {
	err := dp.Flush()
	if err != nil {
		return fmt.Errorf("flush: %v", err)
	}

	dp.dataset = nil
	dp.isLoaded = false

	return nil
}

func (dp *datasetPartition) Flush() error {
	if !dp.isModified {
		return nil
	}

	ctx := context.TODO()

	iter, err := dp.dataset.NewStatementIterator(ctx)
	if err != nil {
		return fmt.Errorf("new iterator: %v", err)
	}

	filePath := dp.getFilePath()

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %v", err)
	}

	fh, err := os.OpenFile(dp.getFilePath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open: %v", err)
	}

	defer fh.Close()

	w, err := dp.dpm.documentEncoding.NewDatasetEncoder(fh)
	if err != nil {
		return fmt.Errorf("new writer: %v", err)
	}

	defer w.Close()

	if wResource, ok := w.(interface {
		PutResource(ctx context.Context, r rdfdescription.Resource) error
	}); ok {
		builder := rdfdescription.NewResourceListBuilder()

		for iter.Next() {
			builder.AddTriple(iter.GetTriple())
		}

		if err := iter.Err(); err != nil {
			return fmt.Errorf("read: %v", err)
		}

		for _, r := range builder.GetResources() {
			err := wResource.PutResource(ctx, r)
			if err != nil {
				return fmt.Errorf("write: %v", err)
			}
		}
	} else {
		for iter.Next() {
			err := w.PutGraphTriple(ctx, iter.GetGraphName(), iter.GetTriple())
			if err != nil {
				return fmt.Errorf("write: %v", err)
			}
		}

		if err := iter.Err(); err != nil {
			return fmt.Errorf("read: %v", err)
		}
	}

	dp.isModified = false

	return nil
}

func (dp *datasetPartition) requireLoad(ctx context.Context) error {
	if dp.isLoaded {
		return nil
	}

	fh, err := os.OpenFile(dp.getFilePath(), os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dp.dataset = inmemory.NewDataset()
			dp.isLoaded = true

			return nil
		}

		return fmt.Errorf("open: %v", err)
	}

	defer fh.Close()

	r, err := dp.dpm.documentEncoding.NewDatasetDecoder(fh)
	if err != nil {
		return fmt.Errorf("new reader: %v", err)
	}

	dp.dataset = inmemory.NewDataset()

	for r.Next() {
		err := dp.dataset.ImportStatement(ctx, r.GetStatement())
		if err != nil {
			return fmt.Errorf("import: %v", err)
		}
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("read: %v", err)
	}

	dp.isLoaded = true

	return nil
}

func (dp *datasetPartition) NewGraphIterator(ctx context.Context) (rdfio.GraphIterator, error) {
	err := dp.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return dp.dataset.NewGraphIterator(ctx)
}

func (dp *datasetPartition) NewGraphNameIterator(ctx context.Context) (rdfio.GraphNameIterator, error) {
	err := dp.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return dp.dataset.NewGraphNameIterator(ctx)
}

func (dp *datasetPartition) NewNodeIterator(ctx context.Context) (rdfio.NodeIterator, error) {
	err := dp.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return dp.dataset.NewNodeIterator(ctx)
}

func (dp *datasetPartition) NewStatementIterator(ctx context.Context, matchers ...rdfio.StatementMatcher) (rdfio.DatasetStatementIterator, error) {
	err := dp.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return dp.dataset.NewStatementIterator(ctx, matchers...)
}

func (dp *datasetPartition) GetGraph(ctx context.Context, graphName rdf.GraphNameValue) rdfio.Graph {
	return &datasetPartitionGraph{
		dp:        dp,
		graphName: graphName,
	}
}

func (dp *datasetPartition) GetStatement(ctx context.Context, triple rdf.Triple) (rdfio.Statement, error) {
	err := dp.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return dp.dataset.GetStatement(ctx, triple)
}

func (dp *datasetPartition) GetGraphStatement(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) (rdfio.Statement, error) {
	err := dp.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return dp.dataset.GetGraphStatement(ctx, graphName, triple)
}

func (dp *datasetPartition) PutTriple(ctx context.Context, triple rdf.Triple) error {
	err := dp.requireLoad(ctx)
	if err != nil {
		return err
	}

	dp.isModified = true

	return dp.dataset.PutTriple(ctx, triple)
}

func (dp *datasetPartition) DeleteTriple(ctx context.Context, triple rdf.Triple) error {
	err := dp.requireLoad(ctx)
	if err != nil {
		return err
	}

	dp.isModified = true

	return dp.dataset.DeleteTriple(ctx, triple)
}

func (dp *datasetPartition) DeleteGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	err := dp.requireLoad(ctx)
	if err != nil {
		return err
	}

	dp.isModified = true

	return dp.dataset.DeleteGraphTriple(ctx, graphName, triple)
}

func (dp *datasetPartition) PutGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	err := dp.requireLoad(ctx)
	if err != nil {
		return err
	}

	dp.isModified = true

	return dp.dataset.PutGraphTriple(ctx, graphName, triple)
}
