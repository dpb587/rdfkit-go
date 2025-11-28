package discard

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionio"
	"github.com/dpb587/rdfkit-go/rdfio"
)

var Encoder = &encoder{}

type encoder struct{}

var _ encoding.GraphEncoder = &encoder{}
var _ rdfdescriptionio.GraphEncoder = &encoder{}

func (w *encoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{}
}

func (w *encoder) Close() error {
	return nil
}

func (w *encoder) PutTriple(_ context.Context, _ rdf.Triple) error {
	return nil
}

func (w *encoder) PutStatement(_ context.Context, _ rdfio.Statement) error {
	return nil
}

func (w *encoder) PutResource(_ context.Context, _ rdfdescription.Resource) error {
	return nil
}
