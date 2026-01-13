package encodingtest

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

const DiscardEncoderContentTypeIdentifier encoding.ContentTypeIdentifier = "internal.dev.discard"

var DiscardEncoder = &discardEncoder{}

type discardEncoder struct{}

var _ encoding.QuadsEncoder = &discardEncoder{}
var _ encoding.TriplesEncoder = &discardEncoder{}
var _ rdfdescription.ResourceWriter = &discardEncoder{}
var _ rdfdescription.DatasetResourceWriter = &discardEncoder{}

func (w *discardEncoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return DiscardEncoderContentTypeIdentifier
}

func (w *discardEncoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{}
}

func (w *discardEncoder) Close() error {
	return nil
}

func (w *discardEncoder) AddQuad(_ context.Context, _ rdf.Quad) error {
	return nil
}

func (w *discardEncoder) AddTriple(_ context.Context, _ rdf.Triple) error {
	return nil
}

func (w *discardEncoder) AddResource(_ context.Context, _ rdfdescription.Resource) error {
	return nil
}

func (w *discardEncoder) AddDatasetResource(_ context.Context, _ rdfdescription.DatasetResource) error {
	return nil
}
