package encodingutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type TripleAsQuadDecoder struct {
	encoding.TriplesDecoder
	GraphName rdf.GraphNameValue
}

var _ encoding.QuadsDecoder = TripleAsQuadDecoder{}
var _ encoding.TriplesDecoder = TripleAsQuadDecoder{}

func (d TripleAsQuadDecoder) Quad() rdf.Quad {
	return rdf.Quad{
		Triple:    d.TriplesDecoder.Triple(),
		GraphName: d.GraphName,
	}
}

func (d TripleAsQuadDecoder) Statement() rdf.Statement {
	return d.Quad()
}
