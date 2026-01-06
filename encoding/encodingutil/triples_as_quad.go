package encodingutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type TripleAsQuadDecoder struct {
	u            encoding.TriplesDecoder
	uTextOffsets encoding.StatementTextOffsetsProvider
	graphName    rdf.GraphNameValue
}

var _ encoding.QuadsDecoder = TripleAsQuadDecoder{}
var _ encoding.StatementTextOffsetsProvider = TripleAsQuadDecoder{}

func NewTripleAsQuadDecoder(u encoding.TriplesDecoder, graphName rdf.GraphNameValue) TripleAsQuadDecoder {
	d := TripleAsQuadDecoder{
		u:         u,
		graphName: graphName,
	}

	d.uTextOffsets, _ = u.(encoding.StatementTextOffsetsProvider)

	return d
}

func (d TripleAsQuadDecoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return d.u.GetContentTypeIdentifier()
}

func (d TripleAsQuadDecoder) Close() error {
	return d.u.Close()
}

func (d TripleAsQuadDecoder) Err() error {
	return d.u.Err()
}

func (d TripleAsQuadDecoder) Next() bool {
	return d.u.Next()
}

func (d TripleAsQuadDecoder) Quad() rdf.Quad {
	return rdf.Quad{
		Triple:    d.u.Triple(),
		GraphName: d.graphName,
	}
}

func (d TripleAsQuadDecoder) Statement() rdf.Statement {
	return d.Quad()
}

func (d TripleAsQuadDecoder) StatementTextOffsets() encoding.StatementTextOffsets {
	if d.uTextOffsets == nil {
		return nil
	}

	return d.uTextOffsets.StatementTextOffsets()
}
