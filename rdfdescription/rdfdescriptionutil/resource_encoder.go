package rdfdescriptionutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type ResourceEncoder interface {
	encoding.Encoder
	rdfdescription.ResourceWriter
}
