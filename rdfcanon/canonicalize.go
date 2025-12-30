package rdfcanon

import (
	"maps"
	"slices"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type CanonicalizeOption interface {
	apply(s *CanonicalizeConfig)
	newCanonicalizer(input rdf.QuadIterator) (*canonicalizer, error)
}

//

type canonicalizer struct {
	canonicalizationState *canonicalizationState
	buildCanonicalQuad    bool
	input                 rdf.QuadIterator
}

//

func Canonicalize(input rdf.QuadIterator, options ...CanonicalizeOption) (*Canonicalization, error) {
	c := CanonicalizeConfig{}

	for _, opt := range options {
		opt.apply(&c)
	}

	cc, err := c.newCanonicalizer(input)
	if err != nil {
		return nil, err
	}

	return algorithmCanonicalization{
		canonicalizationState: cc.canonicalizationState,
		buildCanonicalQuad:    cc.buildCanonicalQuad,
		input:                 cc.input,
	}.Call()
}

//

type identifierIssuer struct {
	stringer         blanknodeutil.Stringer
	knownIdentifiers map[rdf.BlankNodeIdentifier]string
	issuedOrder      []rdf.BlankNodeIdentifier
}

func (i *identifierIssuer) GetBlankNodeIdentifierIfKnown(bn rdf.BlankNode) (string, bool) {
	id, ok := i.knownIdentifiers[bn.GetBlankNodeIdentifier()]

	return id, ok
}

func (i *identifierIssuer) GetBlankNodeIdentifier(bn rdf.BlankNode) string {
	bnId := bn.GetBlankNodeIdentifier()

	id := i.stringer.GetBlankNodeIdentifier(bn)

	if _, ok := i.knownIdentifiers[bnId]; !ok {
		i.issuedOrder = append(i.issuedOrder, bnId)
		i.knownIdentifiers[bnId] = id
	}

	return id
}

func (i *identifierIssuer) Clone() identifierIssuer {
	return identifierIssuer{
		stringer:         i.stringer,
		knownIdentifiers: maps.Clone(i.knownIdentifiers),
		issuedOrder:      slices.Clone(i.issuedOrder),
	}
}

type canonicalizationQuad struct {
	Original      rdf.Quad
	OriginalIndex int64

	SubjectEncoded             string
	SubjectBlankNodeIdentifier rdf.BlankNodeIdentifier

	PredicateEncoded string

	ObjectEncoded             string
	ObjectBlankNodeIdentifier rdf.BlankNodeIdentifier

	GraphNameEncoded             string
	GraphNameBlankNodeIdentifier rdf.BlankNodeIdentifier
}
