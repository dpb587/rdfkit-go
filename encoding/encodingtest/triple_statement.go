package encodingtest

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type TripleStatement struct {
	Triple      rdf.Triple
	TextOffsets encoding.StatementTextOffsets
}

//

type TripleStatementList []TripleStatement

func (tsl TripleStatementList) AsTriples() rdf.TripleList {
	triples := make(rdf.TripleList, 0, len(tsl))

	for _, ts := range tsl {
		triples = append(triples, ts.Triple)
	}

	return triples
}
