package rdf

import "sync/atomic"

// BlankNodeFactory allocates a new blank node.
type BlankNodeFactory interface {
	NewBlankNode() BlankNode
}

//

var defaultBlankNodeFactory = &builtinBlankNodeFactory{
	a: &atomic.Int64{},
}

// DefaultBlankNodeFactory is the default blank node factory used by [NewBlankNode].
var DefaultBlankNodeFactory BlankNodeFactory = defaultBlankNodeFactory

//

type builtinBlankNodeFactory struct {
	a *atomic.Int64
}

func (g *builtinBlankNodeFactory) NewBlankNode() BlankNode {
	return blankNode{
		identifier: builtinBlankNodeIdentifier{
			i: g.a.Add(1),
		},
	}
}

type builtinBlankNodeIdentifier struct {
	UnimplementedBlankNodeIdentifier

	i int64
}
