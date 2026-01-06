package rdf

import "sync/atomic"

// BlankNodeFactory allocates a new, unique blank node.
type BlankNodeFactory interface {
	NewBlankNode() BlankNode
}

//

func NewBlankNode() BlankNode {
	return DefaultBlankNodeFactory.NewBlankNode()
}

//

// DefaultBlankNodeFactory is the default blank node factory used by [NewBlankNode].
var DefaultBlankNodeFactory BlankNodeFactory = &unscopedBlankNodeFactory{
	a: &atomic.Int64{},
}

//

type blankNodeIdentifier struct {
	scope *blankNodeFactory
	value int64
}

var _ BlankNodeIdentifier = blankNodeIdentifier{}

func (bni blankNodeIdentifier) EqualsBlankNodeIdentifier(other BlankNodeIdentifier) bool {
	otherT, ok := other.(blankNodeIdentifier)
	if !ok {
		return false
	}

	return otherT.scope == bni.scope && otherT.value == bni.value
}

//

type blankNodeFactory struct {
	a *atomic.Int64
}

var _ BlankNodeFactory = &blankNodeFactory{}

func NewBlankNodeFactory() BlankNodeFactory {
	return &blankNodeFactory{
		a: &atomic.Int64{},
	}
}

func (g *blankNodeFactory) NewBlankNode() BlankNode {
	return BlankNode{
		Identifier: blankNodeIdentifier{
			scope: g,
			value: g.a.Add(1),
		},
	}
}
