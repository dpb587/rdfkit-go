package nquads

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type statement struct {
	graphName rdf.GraphNameValue
	triple    rdf.Triple
	offsets   encoding.StatementTextOffsets
}

var _ rdfio.Statement = &statement{}
var _ encoding.DecoderTextOffsetsStatement = &statement{}

func (t *statement) GetGraphName() rdf.GraphNameValue {
	if t.graphName == nil {
		return rdf.DefaultGraph
	}

	return t.graphName
}

func (t *statement) GetTriple() rdf.Triple {
	return t.triple
}

func (t *statement) GetDecoderTextOffsets() encoding.StatementTextOffsets {
	return t.offsets
}
