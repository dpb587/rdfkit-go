package triples

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

func GraphAddAll(ctx context.Context, d GraphWriter, triples ...rdf.Triple) error {
	for _, t := range triples {
		if err := d.AddTriple(ctx, t); err != nil {
			return err
		}
	}

	return nil
}

func GraphDeleteAll(ctx context.Context, d Graph, triples ...rdf.Triple) error {
	for _, t := range triples {
		if err := d.DeleteTriple(ctx, t); err != nil {
			return err
		}
	}

	return nil
}
