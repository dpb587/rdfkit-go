package rdfdescriptionio

import (
	"context"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type DefaultGraphDatasetEncoder struct {
	GraphEncoder
}

var _ DatasetEncoder = &DefaultGraphDatasetEncoder{}

func (w DefaultGraphDatasetEncoder) PutGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	if graphName != rdf.DefaultGraph {
		return fmt.Errorf("invalid graph name: %v", graphName)
	}

	return w.GraphEncoder.PutTriple(ctx, triple)
}

func (w DefaultGraphDatasetEncoder) PutGraphResource(ctx context.Context, graphName rdf.GraphNameValue, r rdfdescription.Resource) error {
	if graphName != rdf.DefaultGraph {
		return fmt.Errorf("invalid graph name: %v", graphName)
	}

	return w.GraphEncoder.PutResource(ctx, r)
}

//

type DefaultGraphDatasetFactory struct {
	encoding.GraphFactory
}

var _ encoding.DatasetFactory = &DefaultGraphDatasetFactory{}

func (e DefaultGraphDatasetFactory) NewDatasetEncoder(w io.Writer) (encoding.DatasetEncoder, error) {
	gw, err := e.NewGraphEncoder(w)
	if err != nil {
		return nil, err
	}

	gwT, ok := gw.(GraphEncoder)
	if !ok {
		return nil, fmt.Errorf("invalid graph writer: %T", gw)
	}

	return DefaultGraphDatasetEncoder{gwT}, nil
}

func (e DefaultGraphDatasetFactory) NewDatasetDecoder(r io.Reader) (encoding.DatasetDecoder, error) {
	gi, err := e.NewGraphDecoder(r)
	if err != nil {
		return nil, err
	}

	return encodingutil.DatasetGraphStatementIterator{
		GraphStatementIterator: gi,
	}, nil
}

func (e DefaultGraphDatasetFactory) GetDatasetEncoderContentMetadata() encoding.ContentMetadata {
	return e.GetGraphEncoderContentMetadata()
}
