package jsonld

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type statement struct {
	graphName         rdf.GraphNameValue
	triple            rdf.Triple
	offsets           encoding.StatementTextOffsets
	containerResource encoding.ContainerResource
}

var _ rdfio.Statement = &statement{}
var _ encoding.DecoderTextOffsetsStatement = &statement{}
var _ encoding.ContainerProvider = &statement{}

func (t *statement) GetGraphName() rdf.GraphNameValue {
	return t.graphName
}

func (t *statement) GetTriple() rdf.Triple {
	return t.triple
}

func (t *statement) GetDecoderTextOffsets() encoding.StatementTextOffsets {
	return t.offsets
}

func (tb *statement) GetEncodingContainer() (encoding.Container, bool) {
	if tb.containerResource != nil {
		return encoding.Container{
			Resource: tb.containerResource,
		}, true
	}

	return encoding.Container{}, false
}
