package blanknodeutil

import (
	"sync/atomic"

	"github.com/dpb587/rdfkit-go/rdf"
)

type builtinFactory struct {
	a *atomic.Int64
}

func NewFactory() rdf.BlankNodeFactory {
	return &builtinFactory{
		a: &atomic.Int64{},
	}
}

func (g *builtinFactory) NewBlankNode() rdf.BlankNode {
	return rdf.NewBlankNodeWithIdentifier(
		builtinFactoryIdentifier{
			g: g,
			i: g.a.Add(1),
		},
	)
}

type builtinFactoryIdentifier struct {
	rdf.UnimplementedBlankNodeIdentifier

	g *builtinFactory
	i int64
}
