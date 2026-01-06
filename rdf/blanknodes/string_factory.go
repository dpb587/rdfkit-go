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

type stringIdentifier struct {
	scope *stringFactory
	value string
}

var _ rdf.BlankNodeIdentifier = stringIdentifier{}

func (bni stringIdentifier) EqualsBlankNodeIdentifier(other rdf.BlankNodeIdentifier) bool {
	otherT, ok := other.(stringIdentifier)
	if !ok {
		return false
	}

	return otherT.scope == bni.scope && otherT.value == bni.value
}

//

type stringFactory struct {
	anon rdf.BlankNodeFactory
}

var _ StringFactory = &stringFactory{}
var _ StringProviderProvider = &stringFactory{}

func NewStringFactory() StringFactory {
	return &stringFactory{
		anon: rdf.NewBlankNodeFactory(),
	}
}

func (bnf *stringFactory) NewBlankNode() rdf.BlankNode {
	return bnf.anon.NewBlankNode()
}

func (bnf *stringFactory) NewStringBlankNode(identifier string) rdf.BlankNode {
	if len(identifier) == 0 {
		return bnf.anon.NewBlankNode()
	}

	return rdf.BlankNode{
		Identifier: stringIdentifier{
			scope: bnf,
			value: identifier,
		},
	}
}

func (bnf *stringFactory) GetStringProvider(fallback StringProvider) StringProvider {
	return stringIdentifierProvider{
		scope:    bnf,
		fallback: fallback,
	}
}

//

type stringIdentifierProvider struct {
	scope    *stringFactory
	fallback StringProvider
}

var _ StringProvider = stringIdentifierProvider{}

func (sp stringIdentifierProvider) GetBlankNodeString(bn rdf.BlankNode) string {
	if bnT, ok := bn.Identifier.(stringIdentifier); ok && bnT.scope == sp.scope {
		return bnT.value
	}

	return sp.fallback.GetBlankNodeString(bn)
}
