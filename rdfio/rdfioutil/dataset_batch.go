package rdfioutil

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

func DatasetPutTriples(ctx context.Context, w rdfio.DatasetWriter, g rdf.GraphNameValue, sl rdf.TripleList) error {
	for _, binding := range sl {
		err := w.PutGraphTriple(ctx, g, binding)
		if err != nil {
			return err
		}
	}

	return nil
}
