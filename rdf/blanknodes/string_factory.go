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
