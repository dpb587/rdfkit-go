package rdfioutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type Statement struct {
	GraphName   rdf.GraphNameValue
	Triple      rdf.Triple
	TextOffsets encoding.StatementTextOffsets
	Baggage     map[any]any
}

var _ rdfio.Statement = &Statement{}
var _ encoding.DecoderTextOffsetsStatement = &Statement{}

func (s Statement) GetGraphName() rdf.GraphNameValue {
	if s.GraphName == nil {
		return rdf.DefaultGraph
	}

	return s.GraphName
}

func (s Statement) GetTriple() rdf.Triple {
	return s.Triple
}

func (s Statement) GetDecoderTextOffsets() encoding.StatementTextOffsets {
	return s.TextOffsets
}
