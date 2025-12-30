package rdfcanon

import (
	"encoding/hex"
	"maps"
	"slices"
	"strings"

	"github.com/cespare/permute/v2"
	"github.com/dpb587/rdfkit-go/rdf"
)

type algorithmHashNDegreeQuads struct {
	canonicalizationState *canonicalizationState
	identifier            rdf.BlankNodeIdentifier
	issuer                identifierIssuer

	// not in spec

	maxRecursionDepth int
}

type algorithmHashNDegreeQuadsResult struct {
	hash             string
	identifierIssuer identifierIssuer
}

func (a algorithmHashNDegreeQuads) Call() (algorithmHashNDegreeQuadsResult, error) {
	if a.maxRecursionDepth < 0 {
		return algorithmHashNDegreeQuadsResult{}, ErrMaxRecursionDepthReached
	}

	// [spec // 4.8.3 // 1] Create a new map Hn for relating hashes to related blank nodes.

	h := map[string][]rdf.BlankNodeIdentifier{}

	// [spec // 4.8.3 // 2] Get a reference, quads, to the list of quads from the map entry for identifier in the blank
	// node to quads map.

	quads := a.canonicalizationState.blankNodeToQuads[a.identifier]

	// [spec // 4.8.3 // 3] For each quad in quads:

	for _, quad := range quads {

		// [spec // 4.8.3 // 3.1] For each component in quad, where component is the subject, object, or graph name, and it
		// is a blank node that is not identified by identifier:

		eachComponent := func(component rdf.BlankNodeIdentifier, position string) {

			// [spec // 4.8.3 // 3.1.1] Set hash to the result of the Hash Related Blank Node algorithm, passing the blank
			// node identifier for component as related, quad, issuer, and position as either s, o, or g based on whether
			// component is a subject, object, graph name, respectively.

			hash := algorithmHashRelatedBlankNode{
				related:  component,
				quad:     quad,
				issuer:   a.issuer,
				position: position,
				// not explicitly mentioned but required
				canonicalizationState: a.canonicalizationState,
			}.Call()

			// [spec // 4.8.3 // 3.1.2] Add a mapping of hash to the blank node identifier for component to Hn, adding an
			// entry as necessary.

			h[hash] = append(h[hash], component)

		}

		if quad.SubjectBlankNodeIdentifier != nil && quad.SubjectBlankNodeIdentifier != a.identifier {
			eachComponent(quad.SubjectBlankNodeIdentifier, "s")
		}

		if quad.ObjectBlankNodeIdentifier != nil && quad.ObjectBlankNodeIdentifier != a.identifier {
			eachComponent(quad.ObjectBlankNodeIdentifier, "o")
		}

		if quad.GraphNameBlankNodeIdentifier != nil && quad.GraphNameBlankNodeIdentifier != a.identifier {
			eachComponent(quad.GraphNameBlankNodeIdentifier, "g")
		}
	}

	// [spec // 4.8.3 // 4] Create an empty string, data to hash.

	var dataToHash string

	// [spec // 4.8.3 // 5] For each related hash to blank node list mapping in Hn, code point ordered by related hash:

	orderedRelatedHashes := slices.Collect(maps.Keys(h))
	slices.SortFunc(orderedRelatedHashes, strings.Compare)

	for _, relatedHash := range orderedRelatedHashes {

		blankNodeList := h[relatedHash]

		// [spec // 4.8.3 // 5.1] Append the related hash to the data to hash.

		dataToHash += relatedHash

		// [spec // 4.8.3 // 5.2] Create a string chosen path.

		var chosenPath string

		// [spec // 4.8.3 // 5.3] Create an unset chosen issuer variable.

		var chosenIssuer identifierIssuer

		// [spec // 4.8.3 // 5.4] For each permutation p of blank node list:

		blankNodeListPermutations := permute.Slice(blankNodeList)

		for blankNodeListPermutations.Permute() {
			a.canonicalizationState.maxPermutations--

			if a.canonicalizationState.maxPermutations < 0 {
				return algorithmHashNDegreeQuadsResult{}, ErrMaxIterationsReached
			}

			p := blankNodeList // permute modifies in place; aliased for consistency with spec

			// [spec // 4.8.3 // 5.4.1] Create a copy of issuer, issuer copy.

			issuerCopy := a.issuer.Clone()

			// [spec // 4.8.3 // 5.4.2] Create a string path.

			var path string

			// [spec // 4.8.3 // 5.4.3] Create a recursion list, to store blank node identifiers that must be recursively
			// processed by this algorithm.

			var recursionList []rdf.BlankNodeIdentifier

			// [spec // 4.8.3 // 5.4.4] For each related in p:

			for _, related := range p {

				// [spec // 4.8.3 // 5.4.4.1] If a canonical identifier has been issued for related by canonical issuer, append
				// the string _:, followed by the canonical identifier for related, to path.

				if id, ok := a.canonicalizationState.canonicalIssuer.GetBlankNodeIdentifierIfKnown(rdf.NewBlankNodeWithIdentifier(related)); ok {
					path += "_:" + id
				} else {

					// [spec // 4.8.3 // 5.4.4.2] Otherwise:
					// [spec // 4.8.3 // 5.4.4.2.1] If issuer copy has not issued an identifier for related, append related to
					// recursion list.

					if _, ok := issuerCopy.GetBlankNodeIdentifierIfKnown(rdf.NewBlankNodeWithIdentifier(related)); !ok {
						recursionList = append(recursionList, related)
					}

					// [spec // 4.8.3 // 5.4.4.2.2] Use the Issue Identifier algorithm, passing issuer copy and the related, and
					// append the string _:, followed by the result, to path.

					path += "_:" + issuerCopy.GetBlankNodeIdentifier(rdf.NewBlankNodeWithIdentifier(related))

				}

				// [spec // 4.8.3 // 5.4.4.3] If chosen path is not empty and the length of path is greater than or equal to the
				// length of chosen path and path is greater than chosen path when considering code point order, then skip to
				// the next permutation p.

				if len(chosenPath) > 0 && (len(path) >= len(chosenPath)) && (strings.Compare(path, chosenPath) > 0) {
					goto PERMUTATION_NEXT
				}

			}

			// [spec // 4.8.3 // 5.4.5] For each related in recursion list:

			for _, related := range recursionList {

				// [spec // 4.8.3 // 5.4.5.1] Set result to the result of recursively executing the Hash N-Degree Quads
				// algorithm, passing the canonicalization state, related for identifier, and issuer copy for path identifier
				// issuer.

				result, err := algorithmHashNDegreeQuads{
					canonicalizationState: a.canonicalizationState,
					identifier:            related,
					issuer:                issuerCopy,
					maxRecursionDepth:     a.maxRecursionDepth - 1,
				}.Call()
				if err != nil {
					return algorithmHashNDegreeQuadsResult{}, err
				}

				// [spec // 4.8.3 // 5.4.5.2] Use the Issue Identifier algorithm, passing issuer copy and related; append the
				// string _:, followed by the result, to path.

				path += "_:" + issuerCopy.GetBlankNodeIdentifier(rdf.NewBlankNodeWithIdentifier(related))

				// [spec // 4.8.3 // 5.4.5.3] Append <, the hash in result, and > to path.

				path += "<" + result.hash + ">"

				// [spec // 4.8.3 // 5.4.5.4] Set issuer copy to the identifier issuer in result.

				issuerCopy = result.identifierIssuer

				// [spec // 4.8.3 // 5.4.5.5] If chosen path is not empty and the length of path is greater than or equal to the
				// length of chosen path and path is greater than chosen path when considering code point order, then skip to
				// the next p.

				if len(chosenPath) > 0 && (len(path) >= len(chosenPath)) && (strings.Compare(path, chosenPath) > 0) {
					goto PERMUTATION_NEXT
				}

			}

			// [spec // 4.8.3 // 5.4.6] If chosen path is empty or path is less than chosen path when considering code point
			// order, set chosen path to path and chosen issuer to issuer copy.

			if len(chosenPath) == 0 || (strings.Compare(path, chosenPath) < 0) {
				chosenPath = path
				chosenIssuer = issuerCopy
			}

		PERMUTATION_NEXT:
		}

		// [spec // 4.8.3 // 5.5] Append chosen path to data to hash.

		dataToHash += chosenPath

		// [spec // 4.8.3 // 5.6] Replace issuer, by reference, with chosen issuer.

		a.issuer = chosenIssuer
	}

	// [spec // 4.8.3 // 6] Return issuer and the hash that results from passing data to hash through the hash algorithm.

	hh := a.canonicalizationState.hashNewer()
	hh.Write([]byte(dataToHash))

	return algorithmHashNDegreeQuadsResult{
		hash:             hex.EncodeToString(hh.Sum(nil)),
		identifierIssuer: a.issuer,
	}, nil
}
