package rdfcanon

import (
	"context"
	"crypto/sha256"
	"hash"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type CanonicalizeConfig struct {
	hashNewer          func() hash.Hash
	bnStringProvider   blanknodes.StringProvider
	buildCanonicalQuad *bool
}

var _ CanonicalizeOption = CanonicalizeConfig{}

func (c CanonicalizeConfig) SetHashFunc(v func() hash.Hash) CanonicalizeConfig {
	c.hashNewer = v

	return c
}

func (c CanonicalizeConfig) SetBlankNodeStringProvider(v blanknodes.StringProvider) CanonicalizeConfig {
	c.bnStringProvider = v

	return c
}

func (c CanonicalizeConfig) SetBuildCanonicalQuad(v bool) CanonicalizeConfig {
	c.buildCanonicalQuad = &v

	return c
}

func (c CanonicalizeConfig) apply(s *CanonicalizeConfig) {
	if c.hashNewer != nil {
		s.hashNewer = c.hashNewer
	}

	if c.bnStringProvider != nil {
		s.bnStringProvider = c.bnStringProvider
	}

	if c.buildCanonicalQuad != nil {
		s.buildCanonicalQuad = c.buildCanonicalQuad
	}
}

func (c CanonicalizeConfig) newCanonicalizer(ctx context.Context, input rdf.QuadIterator) (*canonicalizer, error) {
	cc := &canonicalizer{
		canonicalizationState: &canonicalizationState{
			ctx:              ctx,
			hashNewer:        c.hashNewer,
			blankNodeToQuads: map[rdf.BlankNodeIdentifier][]*canonicalizationQuad{},
			hashToBlankNodes: map[string][]rdf.BlankNodeIdentifier{},
			canonicalIssuer: identifierIssuer{
				knownIdentifiers: map[rdf.BlankNodeIdentifier]string{},
			},
			// these are arbitrary; probably should be configurable?
			maxPermutations:   4096,
			maxRecursionDepth: 512,
		},
		input: input,
	}

	if cc.canonicalizationState.hashNewer == nil {
		cc.canonicalizationState.hashNewer = sha256.New
	}

	if c.bnStringProvider != nil {
		cc.canonicalizationState.canonicalIssuer.stringer = c.bnStringProvider
	} else {
		cc.canonicalizationState.canonicalIssuer.stringer = blanknodes.NewInt64StringProvider("c14n%d")
	}

	if c.buildCanonicalQuad != nil {
		cc.buildCanonicalQuad = *c.buildCanonicalQuad
	}

	return cc, nil
}
