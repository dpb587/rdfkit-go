package rdfdescription

import (
	"context"
)

type ResourceWriter interface {
	AddResource(ctx context.Context, r Resource) error
}

type DatasetResourceWriter interface {
	AddDatasetResource(ctx context.Context, dr DatasetResource) error
}
