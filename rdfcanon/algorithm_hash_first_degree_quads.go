package rdfcanon

import (
	"encoding/hex"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/rdf"
)

type algorithmHashFirstDegreeQuads struct {
	canonicalizationState *canonicalizationState
	input                 rdf.BlankNodeIdentifier
}

func (a algorithmHashFirstDegreeQuads) Call() string {

	// [spec // 4.6.3 // 1] Initialize nquads to an empty list. It will be used to store quads in canonical n-quads form.

	var nquadsList []string

	// [spec // 4.6.3 // 2] Get the list of quads quads from the map entry for reference blank node identifier in the
	// blank node to quads map.

	quads := a.canonicalizationState.blankNodeToQuads[a.input]

	// [spec // 4.6.3 // 3] For each quad quad in quads:

	for _, quad := range quads {

		// [spec // 4.6.3 // 3.1] Serialize the quad in canonical n-quads form with the following special rule:
		// [spec // 4.6.3 // 3.1.1] If any component in quad is an blank node, then serialize it using a special identifier
		// as follows:
		// [spec // 4.6.3 // 3.1.1.1] If the blank node's existing blank node identifier matches the reference blank node
		// identifier then use the blank node identifier a, otherwise, use the blank node identifier z.

		s := &strings.Builder{}

		if len(quad.SubjectEncoded) > 0 {
			s.WriteString(quad.SubjectEncoded)
		} else if quad.SubjectBlankNodeIdentifier == a.input {
			s.WriteString("_:a")
		} else {
			s.WriteString("_:z")
		}

		s.WriteString(" ")

		s.WriteString(quad.PredicateEncoded)
		s.WriteString(" ")

		if len(quad.ObjectEncoded) > 0 {
			s.WriteString(quad.ObjectEncoded)
		} else if quad.ObjectBlankNodeIdentifier == a.input {
			s.WriteString("_:a")
		} else {
			s.WriteString("_:z")
		}

		if len(quad.GraphNameEncoded) > 0 {
			s.WriteString(" ")
			s.WriteString(quad.GraphNameEncoded)
		} else if quad.GraphNameBlankNodeIdentifier == nil {
			// default graph
		} else if quad.GraphNameBlankNodeIdentifier == a.input {
			s.WriteString(" _:a")
		} else {
			s.WriteString(" _:z")
		}

		s.WriteString(" .\n")

		nquadsList = append(nquadsList, s.String())
	}

	// [spec // 4.6.3 // 4] Sort nquads in Unicode code point order.

	slices.SortFunc(nquadsList, strings.Compare)

	// [spec // 4.6.3 // 5] Return the hash that results from passing the sorted and concatenated nquads through the hash
	// algorithm.

	h := a.canonicalizationState.hashNewer()

	for _, nquad := range nquadsList {
		h.Write([]byte(nquad))
	}

	return hex.EncodeToString(h.Sum(nil))
}
