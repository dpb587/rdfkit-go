package rdfcanon

import (
	"encoding/hex"

	"github.com/dpb587/rdfkit-go/rdf"
)

type algorithmHashRelatedBlankNode struct {
	canonicalizationState *canonicalizationState
	related               rdf.BlankNodeIdentifier
	quad                  *canonicalizationQuad
	issuer                identifierIssuer
	position              string
}

func (a algorithmHashRelatedBlankNode) Call() string {

	// [spec // 4.7.3 // 1] Initialize a string input to the value of position.

	input := a.position

	// [spec // 4.7.3 // 2] If position is not g, append <, the value of the predicate in quad, and > to input.

	if a.position != "g" {
		input += "<" + a.quad.PredicateEncoded + ">"
	}

	// [spec // 4.7.3 // 3] If there is a canonical identifier for related, or an identifier issued by issuer, append the
	// string _:, followed by that identifier (using the canonical identifier if present, otherwise the one issued by
	// issuer) to input.

	if id, ok := a.canonicalizationState.canonicalIssuer.GetBlankNodeIdentifierIfKnown(rdf.NewBlankNodeWithIdentifier(a.related)); ok {
		input += "_:" + id
	} else if id, ok := a.issuer.GetBlankNodeIdentifierIfKnown(rdf.NewBlankNodeWithIdentifier(a.related)); ok {
		input += "_:" + id
	} else {

		// [spec // 4.7.3 // 4] Otherwise, append the result of the Hash First Degree Quads algorithm, passing related to input.

		input += algorithmHashFirstDegreeQuads{
			canonicalizationState: a.canonicalizationState,
			input:                 a.related,
		}.Call()

	}

	// [spec // 4.7.3 // 5] Return the hash that results from passing input through the hash algorithm.

	h := a.canonicalizationState.hashNewer()
	h.Write([]byte(input))

	return hex.EncodeToString(h.Sum(nil))
}
