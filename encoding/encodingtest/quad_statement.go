package encodingtest

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type QuadStatement struct {
	Quad        rdf.Quad
	TextOffsets encoding.StatementTextOffsets
}

//

type QuadStatementList []QuadStatement

func (qsl QuadStatementList) AsQuads() rdf.QuadList {
	quads := make(rdf.QuadList, 0, len(qsl))

	for _, qs := range qsl {
		quads = append(quads, qs.Quad)
	}

	return quads
}
