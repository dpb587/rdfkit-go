package fs

import "github.com/dpb587/rdfkit-go/rdf"

type DatasetPartitionFilePather interface {
	GetDatasetPartitionFilePath(key rdf.IRI, graphName rdf.GraphNameValue) (string, error)
}

//

type DatasetPartitionFilePatherFunc func(key rdf.IRI, graphName rdf.GraphNameValue) (string, error)

func (f DatasetPartitionFilePatherFunc) GetDatasetPartitionFilePath(key rdf.IRI, graphName rdf.GraphNameValue) (string, error) {
	return f(key, graphName)
}
