package blanknodes

import "github.com/dpb587/rdfkit-go/rdf"

// StringFactory expands supports for the creation of BlankNode values based on string identifiers.
type StringFactory interface {
	rdf.BlankNodeFactory

	// NewStringBlankNode returns a BlankNode whose equivalence is determined by the identifier string. If identifier is
	// empty, it is equivalent to calling NewBlankNode.
	NewStringBlankNode(identifier string) rdf.BlankNode
}

//

type bnString struct {
	v string
	s *bnStringF
}

var _ rdf.BlankNodeIdentifier = bnString{}

func (bni bnString) EqualsBlankNodeIdentifier(other rdf.BlankNodeIdentifier) bool {
	otherT, ok := other.(bnString)
	if !ok {
		return false
	}

	return otherT.s == bni.s && otherT.v == bni.v
}

//

type bnStringF struct {
	anon rdf.BlankNodeFactory
}

var _ StringFactory = &bnStringF{}
var _ StringProviderProvider = &bnStringF{}

func NewStringFactory() StringFactory {
	return &bnStringF{
		anon: rdf.NewBlankNodeFactory(),
	}
}

func (bnf *bnStringF) NewBlankNode() rdf.BlankNode {
	return bnf.anon.NewBlankNode()
}

func (bnf *bnStringF) NewStringBlankNode(identifier string) rdf.BlankNode {
	if len(identifier) == 0 {
		return bnf.anon.NewBlankNode()
	}

	return rdf.BlankNode{
		Identifier: bnString{
			v: identifier,
			s: bnf,
		},
	}
}

func (bnf *bnStringF) GetStringProvider(fallback StringProvider) StringProvider {
	return stringIdentifierProvider{
		scope:    bnf,
		fallback: fallback,
	}
}

//

type stringIdentifierProvider struct {
	scope    *bnStringF
	fallback StringProvider
}

var _ StringProvider = stringIdentifierProvider{}

func (sp stringIdentifierProvider) GetBlankNodeString(bn rdf.BlankNode) string {
	if bnT, ok := bn.Identifier.(bnString); ok && bnT.s == sp.scope {
		return bnT.v
	}

	return sp.fallback.GetBlankNodeString(bn)
}
