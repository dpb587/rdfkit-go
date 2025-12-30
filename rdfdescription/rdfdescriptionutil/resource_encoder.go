package rdfdescriptionutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type ResourceEncoder interface {
	rdfdescription.ResourceWriter

	Close() error
	GetContentMetadata() encoding.ContentMetadata
}
