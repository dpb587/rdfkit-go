package rdf

import (
	"sync/atomic"
)

// BlankNodeFactory is used to allocate new, globally-unique blank nodes.
type BlankNodeFactory interface {
	NewBlankNode() BlankNode
}

//

// NewBlankNode creates a blank node using [DefaultBlankNodeFactory].
func NewBlankNode() BlankNode {
	return DefaultBlankNodeFactory.NewBlankNode()
}

//

// DefaultBlankNodeFactory is the default blank node factory used by [NewBlankNode].
var DefaultBlankNodeFactory BlankNodeFactory = &defaultBlankNodeFactory{
	a: &atomic.Int64{},
}

//

type bn struct {
	v int64
	s *bnF
}

var _ BlankNodeIdentifier = bn{}

func (bni bn) EqualsBlankNodeIdentifier(other BlankNodeIdentifier) bool {
	otherT, ok := other.(bn)
	if !ok {
		return false
	}

	return otherT.s == bni.s && otherT.v == bni.v
}

//

type bnF struct {
	a *atomic.Int64
}

var _ BlankNodeFactory = &bnF{}

// NewBlankNodeFactory creates a new [BlankNodeFactory].
func NewBlankNodeFactory() BlankNodeFactory {
	return &bnF{
		a: &atomic.Int64{},
	}
}

func (g *bnF) NewBlankNode() BlankNode {
	return BlankNode{
		Identifier: bn{
			v: g.a.Add(1),
			s: g,
		},
	}
}
