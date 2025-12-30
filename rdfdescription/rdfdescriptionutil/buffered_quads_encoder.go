package rdfdescriptionutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type BufferedQuadsEncoder struct {
	ctx        context.Context
	encoder    DatasetResourceEncoder
	preferAnon bool

	builder *rdfdescription.DatasetResourceListBuilder
}

var _ encoding.QuadsEncoder = &BufferedQuadsEncoder{}

func NewBufferedQuadsEncoder(ctx context.Context, encoder DatasetResourceEncoder, preferAnon bool) *BufferedQuadsEncoder {
	return &BufferedQuadsEncoder{
		ctx:        ctx,
		encoder:    encoder,
		preferAnon: preferAnon,
		builder:    rdfdescription.NewDatasetResourceListBuilder(),
	}
}

func (e *BufferedQuadsEncoder) AddQuad(ctx context.Context, quad rdf.Quad) error {
	return e.builder.AddQuad(ctx, quad)
}

func (e *BufferedQuadsEncoder) Close() error {
	err := e.builder.AddToDataset(e.ctx, e.encoder, e.preferAnon)
	if err != nil {
		return err
	}

	return e.encoder.Close()
}

func (e *BufferedQuadsEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.encoder.GetContentMetadata()
}
