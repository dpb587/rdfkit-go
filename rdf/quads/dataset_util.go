package quads

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

func DatasetAddAll(ctx context.Context, d DatasetWriter, quads ...rdf.Quad) error {
	for _, t := range quads {
		if err := d.AddQuad(ctx, t); err != nil {
			return err
		}
	}

	return nil
}

func DatasetDeleteAll(ctx context.Context, d Dataset, quads ...rdf.Quad) error {
	for _, t := range quads {
		if err := d.DeleteQuad(ctx, t); err != nil {
			return err
		}
	}

	return nil
}
