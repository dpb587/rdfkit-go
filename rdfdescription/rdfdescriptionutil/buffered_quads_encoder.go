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
	exportOpts rdfdescription.ExportResourceOptions

	builder *rdfdescription.DatasetResourceListBuilder
}

var _ encoding.QuadsEncoder = &BufferedQuadsEncoder{}

func NewBufferedQuadsEncoder(ctx context.Context, encoder DatasetResourceEncoder, exportOpts rdfdescription.ExportResourceOptions) *BufferedQuadsEncoder {
	return &BufferedQuadsEncoder{
		ctx:        ctx,
		encoder:    encoder,
		exportOpts: exportOpts,
		builder:    rdfdescription.NewDatasetResourceListBuilder(),
	}
}

func (e *BufferedQuadsEncoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return e.encoder.GetContentTypeIdentifier()
}

func (e *BufferedQuadsEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.encoder.GetContentMetadata()
}

func (e *BufferedQuadsEncoder) AddQuad(ctx context.Context, quad rdf.Quad) error {
	return e.builder.AddQuad(ctx, quad)
}

func (e *BufferedQuadsEncoder) Close() error {
	err := e.builder.ToDatasetResourceWriter(e.ctx, e.encoder, e.exportOpts)
	if err != nil {
		return err
	}

	return e.encoder.Close()
}
