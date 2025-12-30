package rdfdescriptionutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type DatasetResourceEncoder interface {
	rdfdescription.DatasetResourceWriter

	Close() error
	GetContentMetadata() encoding.ContentMetadata
}
