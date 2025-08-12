package rdfdescriptionio

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type BufferedDatasetEncoder struct {
	w                 DatasetEncoder
	wClosePropagation bool
	useAnonResource   bool

	builderByGraphName map[rdf.GraphNameValue]*rdfdescription.ResourceListBuilder
}

var _ DatasetEncoder = &BufferedDatasetEncoder{}

func NewBufferedDatasetEncoder(w DatasetEncoder) *BufferedDatasetEncoder {
	return &BufferedDatasetEncoder{
		w:                  w,
		builderByGraphName: map[rdf.GraphNameValue]*rdfdescription.ResourceListBuilder{},
	}
}

func (e *BufferedDatasetEncoder) SetClosePropagation(v bool) {
	e.wClosePropagation = v
}

func (e *BufferedDatasetEncoder) SetAnonResource(v bool) {
	e.useAnonResource = v
}

func (e *BufferedDatasetEncoder) Close() error {
	ctx := context.Background()

	for graphName, builder := range e.builderByGraphName {
		for _, resource := range builder.GetResources() {
			if e.useAnonResource {
				if bn, ok := resource.GetResourceSubject().(rdf.BlankNode); ok && builder.GetBlankNodeReferences(bn) == 0 {
					resource = rdfdescription.AnonResource{
						Statements: resource.GetResourceStatements(),
					}
				}
			}

			if err := e.w.PutGraphResource(ctx, graphName, resource); err != nil {
				return err
			}
		}
	}

	if e.wClosePropagation {
		return e.w.Close()
	}

	return nil
}

func (e *BufferedDatasetEncoder) PutTriple(ctx context.Context, triple rdf.Triple) error {
	if e.builderByGraphName[rdf.DefaultGraph] == nil {
		e.builderByGraphName[rdf.DefaultGraph] = rdfdescription.NewResourceListBuilder()
	}

	e.builderByGraphName[rdf.DefaultGraph].AddTriple(triple)

	return nil
}

func (e *BufferedDatasetEncoder) PutGraphTriple(ctx context.Context, graphName rdf.GraphNameValue, triple rdf.Triple) error {
	if e.builderByGraphName[graphName] == nil {
		e.builderByGraphName[graphName] = rdfdescription.NewResourceListBuilder()
	}

	e.builderByGraphName[graphName].AddTriple(triple)

	return nil
}

func (e *BufferedDatasetEncoder) GetContentMetadata() encoding.ContentMetadata {
	return e.w.GetContentMetadata()
}

func (e *BufferedDatasetEncoder) PutGraphResource(ctx context.Context, graphName rdf.GraphNameValue, r rdfdescription.Resource) error {
	if e.builderByGraphName[graphName] == nil {
		e.builderByGraphName[graphName] = rdfdescription.NewResourceListBuilder()
	}

	b := e.builderByGraphName[graphName]

	for _, triple := range r.AsTriples() {
		b.AddTriple(triple)
	}

	return nil
}
