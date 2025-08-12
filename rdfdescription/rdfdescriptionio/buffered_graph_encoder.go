package rdfdescriptionio

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type BufferedGraphEncoder struct {
	w                 GraphEncoder
	wClosePropagation bool
	useAnonResource   bool

	builder *rdfdescription.ResourceListBuilder
}

var _ GraphEncoder = &BufferedGraphEncoder{}

func NewBufferedGraphEncoder(w GraphEncoder) *BufferedGraphEncoder {
	return &BufferedGraphEncoder{
		w:                 w,
		wClosePropagation: true,
		builder:           rdfdescription.NewResourceListBuilder(),
	}
}

func (e *BufferedGraphEncoder) SetClosePropagation(v bool) {
	e.wClosePropagation = v
}

func (e *BufferedGraphEncoder) SetAnonResource(v bool) {
	e.useAnonResource = v
}

func (e *BufferedGraphEncoder) Close() error {
	ctx := context.Background()

	for _, resource := range e.builder.GetResources() {
		if e.useAnonResource {
			if bn, ok := resource.GetResourceSubject().(rdf.BlankNode); ok && e.builder.GetBlankNodeReferences(bn) == 0 {
				resource = rdfdescription.AnonResource{
					Statements: resource.GetResourceStatements(),
				}
			}
		}

		if err := e.w.PutResource(ctx, resource); err != nil {
			return err
		}
	}

	if e.wClosePropagation {
		return e.w.Close()
	}

	return nil
}

func (e *BufferedGraphEncoder) PutTriple(ctx context.Context, triple rdf.Triple) error {
	e.builder.AddTriple(triple)

	return nil
}

func (e *BufferedGraphEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.w.GetContentMetadata()
}

func (e *BufferedGraphEncoder) PutResource(ctx context.Context, r rdfdescription.Resource) error {
	for _, triple := range r.AsTriples() {
		e.builder.AddTriple(triple)
	}

	return nil
}
