package fs

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

//

type DatasetPartitionManager struct {
	log              *slog.Logger
	dir              string
	filePather       DatasetPartitionFilePather
	documentEncoding encoding.DatasetFactory

	partitions map[string]*datasetPartition
}

func NewDatasetPartitionManager(
	log *slog.Logger,
	dir string,
	filePather DatasetPartitionFilePather,
	documentEncoding encoding.DatasetFactory,
) *DatasetPartitionManager {
	return &DatasetPartitionManager{
		log:              log,
		dir:              dir,
		filePather:       filePather,
		documentEncoding: documentEncoding,
		partitions:       map[string]*datasetPartition{},
	}
}

func (dpm *DatasetPartitionManager) GetPartition(ctx context.Context, key rdf.IRI, graphName rdf.GraphNameValue) (rdfio.Graph, error) {
	fileBase, err := dpm.filePather.GetDatasetPartitionFilePath(key, graphName)
	if err != nil {
		return nil, fmt.Errorf("lookup: %v", err)
	}

	dp, found := dpm.partitions[fileBase]
	if !found {
		dp = &datasetPartition{
			dpm:      dpm,
			fileBase: fileBase,
		}

		dpm.partitions[fileBase] = dp
	}

	return dp.GetGraph(ctx, graphName), nil
}

func (dpm *DatasetPartitionManager) Close() error {
	for _, dp := range dpm.partitions {
		if err := dp.Close(); err != nil {
			return fmt.Errorf("partition[%s]: %v", dp.fileBase, err)
		}

		delete(dpm.partitions, dp.fileBase)
	}

	return nil
}

func (dpm *DatasetPartitionManager) Flush() error {
	for _, dp := range dpm.partitions {
		if err := dp.Flush(); err != nil {
			return fmt.Errorf("partition[%s]: %v", dp.fileBase, err)
		}
	}

	return nil
}
