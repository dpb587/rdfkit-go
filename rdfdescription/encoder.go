package rdfdescription

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
)

type ResourceWriter interface {
	AddResource(ctx context.Context, r Resource) error
}

type DatasetResourceWriter interface {
	AddDatasetResource(ctx context.Context, r Resource, g rdf.GraphNameValue) error
}
