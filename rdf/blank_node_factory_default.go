package rdf

import (
	"sync/atomic"
)

// same as blankNode{Factory,Identifier}, but without a scope field; used for DefaultBlankNodeFactory

type bnDefault struct {
	v int64
}

var _ BlankNodeIdentifier = bnDefault{}

func (bni bnDefault) EqualsBlankNodeIdentifier(other BlankNodeIdentifier) bool {
	otherT, ok := other.(bnDefault)
	if !ok {
		return false
	}

	return otherT.v == bni.v
}

//

type defaultBlankNodeFactory struct {
	a *atomic.Int64
}

func (g *defaultBlankNodeFactory) NewBlankNode() BlankNode {
	return BlankNode{
		Identifier: bnDefault{
			v: g.a.Add(1),
		},
	}
}
