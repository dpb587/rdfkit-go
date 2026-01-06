package rdfcanon

import (
	"bytes"
	"context"
	"fmt"
	"hash"
	"maps"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type canonicalizationState struct {
	ctx       context.Context
	hashNewer func() hash.Hash

	blankNodeToQuads map[rdf.BlankNodeIdentifier][]*canonicalizationQuad
	hashToBlankNodes map[string][]rdf.BlankNodeIdentifier
	canonicalIssuer  identifierIssuer

	// not officially in spec

	allQuads []*canonicalizationQuad

	// "poison" controls

	maxPermutations   int
	maxRecursionDepth int
}

type algorithmCanonicalization struct {
	canonicalizationState *canonicalizationState
	buildCanonicalQuad    bool
	input                 rdf.QuadIterator
}

func (a algorithmCanonicalization) Call() (*Canonicalization, error) {
	// [spec // 4.4.3 // 1] Create the canonicalization state. If the input dataset is an N-Quads document, parse that
	// document into a dataset in the canonicalized dataset, retaining any blank node identifiers used within that
	// document in the input blank node identifier map; otherwise arbitrary identifiers are assigned for each blank node.

	// [spec // 4.4.3 // 2] For every quad Q in input dataset:

	var inputIdx int64

	b := &bytes.Buffer{}

	for a.input.Next() {
		if b := int64(1024); inputIdx%b == b-1 {
			if err := a.canonicalizationState.ctx.Err(); err != nil {
				return nil, err
			}
		}

		q := &canonicalizationQuad{
			Original:      a.input.Quad(),
			OriginalIndex: inputIdx,
		}

		inputIdx++

		// [spec // 4.4.3 // 2.1] For each blank node that is a component of Q, add a reference to Q from the map entry for
		// the blank node identifier identifier in the blank node to quads map, creating a new entry if necessary, using the
		// identifier for the blank node found in the input blank node identifier map.

		switch s := q.Original.Triple.Subject.(type) {
		case rdf.BlankNode:
			q.SubjectBlankNodeIdentifier = s.Identifier
			a.canonicalizationState.blankNodeToQuads[q.SubjectBlankNodeIdentifier] = append(a.canonicalizationState.blankNodeToQuads[q.SubjectBlankNodeIdentifier], q)
		case rdf.IRI:
			b.Reset()
			nquads.WriteIRI(b, s, false)

			q.SubjectEncoded = b.String()
		default:
			panic(fmt.Errorf("invalid subject type: %T", s))
		}

		switch p := q.Original.Triple.Predicate.(type) {
		case rdf.IRI:
			b.Reset()
			nquads.WriteIRI(b, p, false)

			q.PredicateEncoded = b.String()
		default:
			panic(fmt.Errorf("invalid predicate type: %T", p))
		}

		switch o := q.Original.Triple.Object.(type) {
		case rdf.BlankNode:
			q.ObjectBlankNodeIdentifier = o.Identifier
			a.canonicalizationState.blankNodeToQuads[q.ObjectBlankNodeIdentifier] = append(a.canonicalizationState.blankNodeToQuads[q.ObjectBlankNodeIdentifier], q)
		case rdf.IRI:
			b.Reset()
			nquads.WriteIRI(b, o, false)

			q.ObjectEncoded = b.String()
		case rdf.Literal:
			b.Reset()
			nquads.WriteLiteral(b, o, false)

			q.ObjectEncoded = b.String()
		default:
			panic(fmt.Errorf("invalid object type: %T", o))
		}

		if q.Original.GraphName != nil {
			switch g := q.Original.GraphName.(type) {
			case rdf.BlankNode:
				q.GraphNameBlankNodeIdentifier = g.Identifier
				a.canonicalizationState.blankNodeToQuads[q.GraphNameBlankNodeIdentifier] = append(a.canonicalizationState.blankNodeToQuads[q.GraphNameBlankNodeIdentifier], q)
			case rdf.IRI:
				b.Reset()
				nquads.WriteIRI(b, g, false)

				q.GraphNameEncoded = b.String()
			default:
				panic(fmt.Errorf("invalid graph name type: %T", g))
			}
		}

		a.canonicalizationState.allQuads = append(a.canonicalizationState.allQuads, q)
	}

	// [spec // 4.4.3 // 3] For each key n in the blank node to quads map:

	for n := range a.canonicalizationState.blankNodeToQuads {

		// [spec // 4.4.3 // 3.1] Create a hash, hf(n), for n according to the Hash First Degree Quads algorithm.

		hash := algorithmHashFirstDegreeQuads{
			canonicalizationState: a.canonicalizationState,
			input:                 n,
		}.Call()

		// [spec // 4.4.3 // 3.2] Append n to the value associated to hf(n) in hash to blank nodes map, creating a new entry
		// if necessary.

		a.canonicalizationState.hashToBlankNodes[hash] = append(a.canonicalizationState.hashToBlankNodes[hash], n)
	}

	// [spec // 4.4.3 // 4] For each hash to identifier list map entry in hash to blank nodes map, code point ordered by
	// hash:

	orderedHashes := slices.Collect(maps.Keys(a.canonicalizationState.hashToBlankNodes))
	slices.SortFunc(orderedHashes, strings.Compare)

	for _, hash := range orderedHashes {

		// [spec // 4.4.3 // 4.1] If identifier list has more than one entry, continue to the next mapping.

		if len(a.canonicalizationState.hashToBlankNodes[hash]) > 1 {
			continue
		}

		// [spec // 4.4.3 // 4.2] Use the Issue Identifier algorithm, passing canonical issuer and the single blank node
		// identifier, identifier in identifier list to issue a canonical replacement identifier for identifier.

		a.canonicalizationState.canonicalIssuer.GetBlankNodeString(a.canonicalizationState.hashToBlankNodes[hash][0])

		// [spec // 4.4.3 // 4.3] Remove the map entry for hash from the hash to blank nodes map.

		delete(a.canonicalizationState.hashToBlankNodes, hash)
	}

	// [spec // 4.4.3 // 5] For each hash to identifier list map entry in hash to blank nodes map, code point ordered by
	// hash:

	orderedHashes = slices.Collect(maps.Keys(a.canonicalizationState.hashToBlankNodes))
	slices.SortFunc(orderedHashes, strings.Compare)

	for _, hash := range orderedHashes {

		// [spec // 4.4.3 // 5.1] Create hash path list where each item will be a result of running the Hash N-Degree Quads
		// algorithm.

		hashPathList := []algorithmHashNDegreeQuadsResult{}

		// [spec // 4.4.3 // 5.2] For each blank node identifier n in identifier list:

		for _, n := range a.canonicalizationState.hashToBlankNodes[hash] {

			// [spec // 4.4.3 // 5.2.1] If a canonical identifier has already been issued for n, continue to the next blank
			// node identifier.

			if _, ok := a.canonicalizationState.canonicalIssuer.GetBlankNodeStringIfKnown(n); ok {
				continue
			}

			// [spec // 4.4.3 // 5.2.2] Create temporary issuer, an identifier issuer initialized with the prefix b.

			temporaryIssuer := identifierIssuer{
				stringer:         blanknodes.NewInt64StringProvider("b%d"),
				knownIdentifiers: make(map[rdf.BlankNodeIdentifier]string),
			}

			// [spec // 4.4.3 // 5.2.3] Use the Issue Identifier algorithm, passing temporary issuer and n, to issue a new
			// temporary blank node identifier bn to n.

			temporaryIssuer.GetBlankNodeString(n)

			// [spec // 4.4.3 // 5.2.4] Run the Hash N-Degree Quads algorithm, passing the canonicalization state, n for
			// identifier, and temporary issuer, appending the result to the hash path list.

			result, err := algorithmHashNDegreeQuads{
				canonicalizationState: a.canonicalizationState,
				identifier:            n,
				issuer:                temporaryIssuer,
				maxRecursionDepth:     a.canonicalizationState.maxRecursionDepth,
			}.Call()
			if err != nil {
				return nil, err
			}

			hashPathList = append(hashPathList, result)

		}

		// [spec // 4.4.3 // 5.3] For each result in the hash path list, code point ordered by the hash in result:

		slices.SortFunc(hashPathList, func(a, b algorithmHashNDegreeQuadsResult) int {
			return strings.Compare(a.hash, b.hash)
		})

		for _, result := range hashPathList {

			// [spec // 4.4.3 // 5.3.1] For each blank node identifier, existing identifier, that was issued a temporary
			// identifier by identifier issuer in result, issue a canonical identifier, in the same order, using the Issue
			// Identifier algorithm, passing canonical issuer and existing identifier.

			for _, existingIdentifier := range result.identifierIssuer.issuedOrder {
				a.canonicalizationState.canonicalIssuer.GetBlankNodeString(existingIdentifier)
			}

		}

	}

	// [spec // 4.4.3 // 6] Add the issued identifiers map from the canonical issuer to the canonicalized dataset.
	// [dpb] automatic via GetBlankNodeIdentifier

	// [spec // 4.4.3 // 7] Return the serialized canonical form of the canonicalized dataset. Upon request, alternatively
	// (or additionally) return the canonicalized dataset itself, which includes the input blank node identifier map, and
	// issued identifiers map from the canonical issuer.

	cres := &Canonicalization{
		bnStringProvider: a.canonicalizationState.canonicalIssuer.stringer,
		hasCanonicalQuad: a.buildCanonicalQuad,
	}

	for _, quad := range a.canonicalizationState.allQuads {
		s := &strings.Builder{}

		q := quad.Original

		cquad := rdf.Quad{
			Triple:    q.Triple,
			GraphName: q.GraphName,
		}

		if len(quad.SubjectEncoded) > 0 {
			s.WriteString(quad.SubjectEncoded)
		} else {
			s.WriteString("_:")
			s.WriteString(a.canonicalizationState.canonicalIssuer.GetBlankNodeString(quad.SubjectBlankNodeIdentifier))

			cquad.Triple.Subject = rdf.BlankNode{
				Identifier: quad.SubjectBlankNodeIdentifier,
			}
		}

		s.WriteString(" ")

		s.WriteString(quad.PredicateEncoded)
		s.WriteString(" ")

		if len(quad.ObjectEncoded) > 0 {
			s.WriteString(quad.ObjectEncoded)
		} else {
			s.WriteString("_:")
			s.WriteString(a.canonicalizationState.canonicalIssuer.GetBlankNodeString(quad.ObjectBlankNodeIdentifier))

			cquad.Triple.Object = rdf.BlankNode{
				Identifier: quad.ObjectBlankNodeIdentifier,
			}
		}

		if len(quad.GraphNameEncoded) > 0 {
			s.WriteString(" ")
			s.WriteString(quad.GraphNameEncoded)
		} else if quad.GraphNameBlankNodeIdentifier != nil {
			s.WriteString(" _:")
			s.WriteString(a.canonicalizationState.canonicalIssuer.GetBlankNodeString(quad.GraphNameBlankNodeIdentifier))

			cquad.GraphName = rdf.BlankNode{
				Identifier: quad.GraphNameBlankNodeIdentifier,
			}
		}

		s.WriteString(" .\n")

		cqn := canonicalizedQuad{
			originalIndex: quad.OriginalIndex,
			encoded:       []byte(s.String()),
		}

		if cres.hasCanonicalQuad {
			cqn.canonical = &cquad
		}

		cres.nquads = append(cres.nquads, cqn)
	}

	slices.SortFunc(cres.nquads, func(i, j canonicalizedQuad) int {
		return bytes.Compare(i.encoded, j.encoded)
	})

	return cres, nil
}
