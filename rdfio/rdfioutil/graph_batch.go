package rdfioutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

func GraphPutTriples(ctx context.Context, w rdfio.GraphWriter, sl rdf.TripleList) error {
	for _, binding := range sl {
		err := w.PutTriple(ctx, binding)
		if err != nil {
			return err
		}
	}

	return nil
}
