package htmlmicrodata

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type statement struct {
	triple            rdf.Triple
	offsets           encoding.StatementTextOffsets
	containerResource encoding.ContainerResource
}

var _ rdfio.Statement = &statement{}
var _ encoding.DecoderTextOffsetsStatement = &statement{}
var _ encoding.ContainerProvider = &statement{}

func (statement) GetGraphName() rdf.GraphNameValue {
	return rdf.DefaultGraph
}

func (t *statement) GetTriple() rdf.Triple {
	return t.triple
}

func (tb *statement) GetDecoderTextOffsets() encoding.StatementTextOffsets {
	return tb.offsets
}

func (tb *statement) GetEncodingContainer() (encoding.Container, bool) {
	if tb.containerResource != nil {
		return encoding.Container{
			Resource: tb.containerResource,
		}, true
	}

	return encoding.Container{}, false
}
