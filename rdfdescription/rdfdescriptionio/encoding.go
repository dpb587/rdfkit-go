package rdfdescriptionio

import (
	"context"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type DatasetEncoder interface {
	encoding.DatasetEncoder

	PutGraphResource(ctx context.Context, graphName rdf.GraphNameValue, r rdfdescription.Resource) error
}

type GraphEncoder interface {
	encoding.GraphEncoder

	PutResource(ctx context.Context, r rdfdescription.Resource) error
}
