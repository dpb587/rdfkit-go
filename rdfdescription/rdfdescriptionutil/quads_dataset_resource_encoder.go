package rdfdescriptionutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type QuadsDatasetResourceEncoder struct {
	encoder encoding.QuadsEncoder
}

var _ DatasetResourceEncoder = &QuadsDatasetResourceEncoder{}

func NewQuadsDatasetResourceEncoder(encoder encoding.QuadsEncoder) DatasetResourceEncoder {
	return &QuadsDatasetResourceEncoder{
		encoder: encoder,
	}
}

func (e *QuadsDatasetResourceEncoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return e.encoder.GetContentTypeIdentifier()
}

func (e *QuadsDatasetResourceEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.encoder.GetContentMetadata()
}

func (e *QuadsDatasetResourceEncoder) Close() error {
	return e.encoder.Close()
}

func (e *QuadsDatasetResourceEncoder) AddDatasetResource(ctx context.Context, resource rdfdescription.DatasetResource) error {
	for _, statement := range resource.NewQuads() {
		if err := e.encoder.AddQuad(ctx, statement); err != nil {
			return err
		}
	}

	return nil
}
