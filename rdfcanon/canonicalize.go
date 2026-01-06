package rdfcanon

import (
	"context"
	"maps"
	"slices"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type CanonicalizeOption interface {
	apply(s *CanonicalizeConfig)
	newCanonicalizer(ctx context.Context, input rdf.QuadIterator) (*canonicalizer, error)
}

//

type canonicalizer struct {
	canonicalizationState *canonicalizationState
	buildCanonicalQuad    bool
	input                 rdf.QuadIterator
}

//

func Canonicalize(ctx context.Context, input rdf.QuadIterator, options ...CanonicalizeOption) (*Canonicalization, error) {
	c := CanonicalizeConfig{}

	for _, opt := range options {
		opt.apply(&c)
	}

	cc, err := c.newCanonicalizer(ctx, input)
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
	stringer         blanknodes.StringProvider
	knownIdentifiers map[rdf.BlankNodeIdentifier]string
	issuedOrder      []rdf.BlankNodeIdentifier
}

func (i *identifierIssuer) GetBlankNodeStringIfKnown(bni rdf.BlankNodeIdentifier) (string, bool) {
	id, ok := i.knownIdentifiers[bni]

	return id, ok
}

func (i *identifierIssuer) GetBlankNodeString(bni rdf.BlankNodeIdentifier) string {
	id := i.stringer.GetBlankNodeString(rdf.BlankNode{
		Identifier: bni,
	})

	if _, ok := i.knownIdentifiers[bni]; !ok {
		i.issuedOrder = append(i.issuedOrder, bni)
		i.knownIdentifiers[bni] = id
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
