package rdf

import "sync/atomic"

// same as blankNode{Factory,Identifier}, but without a scope field; used for DefaultBlankNodeFactory

type unscopedBlankNodeIdentifier struct {
	value int64
}

var _ BlankNodeIdentifier = unscopedBlankNodeIdentifier{}

func (bni unscopedBlankNodeIdentifier) EqualsBlankNodeIdentifier(other BlankNodeIdentifier) bool {
	otherT, ok := other.(unscopedBlankNodeIdentifier)
	if !ok {
		return false
	}

	return otherT.value == bni.value
}

//

type unscopedBlankNodeFactory struct {
	a *atomic.Int64
}

func (g *unscopedBlankNodeFactory) NewBlankNode() BlankNode {
	return BlankNode{
		Identifier: unscopedBlankNodeIdentifier{
			value: g.a.Add(1),
		},
	}
}
