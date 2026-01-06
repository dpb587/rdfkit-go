package rdfdescriptionutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type BufferedTriplesEncoder struct {
	ctx        context.Context
	encoder    ResourceEncoder
	preferAnon bool

	builder *rdfdescription.ResourceListBuilder
}

var _ encoding.TriplesEncoder = &BufferedTriplesEncoder{}

func NewBufferedTriplesEncoder(ctx context.Context, encoder ResourceEncoder, preferAnon bool) *BufferedTriplesEncoder {
	return &BufferedTriplesEncoder{
		ctx:        ctx,
		encoder:    encoder,
		preferAnon: preferAnon,
		builder:    rdfdescription.NewResourceListBuilder(),
	}
}

func (e *BufferedTriplesEncoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return e.encoder.GetContentTypeIdentifier()
}

func (e *BufferedTriplesEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.encoder.GetContentMetadata()
}

func (e *BufferedTriplesEncoder) AddTriple(ctx context.Context, triple rdf.Triple) error {
	return e.builder.AddTriple(ctx, triple)
}

func (e *BufferedTriplesEncoder) Close() error {
	err := e.builder.AddTo(e.ctx, e.encoder, e.preferAnon)
	if err != nil {
		return err
	}

	return e.encoder.Close()
}
