package encodingutil

import (
	"context"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type DatasetGraphEncoder struct {
	encoding.GraphEncoder
}

var _ encoding.DatasetEncoder = &DatasetGraphEncoder{}

func (w DatasetGraphEncoder) PutGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	if graphName != rdf.DefaultGraph {
		return fmt.Errorf("invalid graph name: %v", graphName)
	}

	return w.GraphEncoder.PutTriple(ctx, triple)
}

//

type DatasetGraphStatementIterator struct {
	rdfio.GraphStatementIterator
}

var _ rdfio.DatasetStatementIterator = &DatasetGraphStatementIterator{}

func (i DatasetGraphStatementIterator) GetGraphName() rdf.GraphNameValue {
	return rdf.DefaultGraph
}

//

type DatasetGraphDocumentEncoding struct {
	encoding.GraphFactory
}

var _ encoding.DatasetFactory = &DatasetGraphDocumentEncoding{}

func (e DatasetGraphDocumentEncoding) NewDatasetEncoder(w io.Writer) (encoding.DatasetEncoder, error) {
	gw, err := e.NewGraphEncoder(w)
	if err != nil {
		return nil, err
	}

	return DatasetGraphEncoder{gw}, nil
}

func (e DatasetGraphDocumentEncoding) NewDatasetDecoder(r io.Reader) (encoding.DatasetDecoder, error) {
	gi, err := e.NewGraphDecoder(r)
	if err != nil {
		return nil, err
	}

	return DatasetGraphStatementIterator{gi}, nil
}

func (e DatasetGraphDocumentEncoding) GetDatasetEncoderContentMetadata() encoding.ContentMetadata {
	return e.GetGraphEncoderContentMetadata()
}
