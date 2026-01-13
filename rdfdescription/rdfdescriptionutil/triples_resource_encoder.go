package rdfdescriptionutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type TriplesResourceEncoder struct {
	encoder encoding.TriplesEncoder
}

var _ ResourceEncoder = &TriplesResourceEncoder{}

func NewTriplesResourceEncoder(encoder encoding.TriplesEncoder) ResourceEncoder {
	return &TriplesResourceEncoder{
		encoder: encoder,
	}
}

func (e *TriplesResourceEncoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return e.encoder.GetContentTypeIdentifier()
}

func (e *TriplesResourceEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.encoder.GetContentMetadata()
}

func (e *TriplesResourceEncoder) Close() error {
	return e.encoder.Close()
}

func (e *TriplesResourceEncoder) AddResource(ctx context.Context, resource rdfdescription.Resource) error {
	for _, statement := range resource.NewTriples() {
		if err := e.encoder.AddTriple(ctx, statement); err != nil {
			return err
		}
	}

	return nil
}
