package blanknodes

import (
	"sync"

	"github.com/dpb587/rdfkit-go/rdf"
)

// Mapper maps blank nodes from one form to another.
type Mapper interface {
	MapBlankNode(t rdf.BlankNode) rdf.BlankNode
}

//

type factoryMapper struct {
	factory rdf.BlankNodeFactory
	mutex   sync.Mutex
	known   map[rdf.BlankNodeIdentifier]rdf.BlankNodeIdentifier
}

func NewFactoryMapper(factory rdf.BlankNodeFactory) Mapper {
	m := &factoryMapper{
		factory: factory,
		known:   map[rdf.BlankNodeIdentifier]rdf.BlankNodeIdentifier{},
	}

	if m.factory == nil {
		m.factory = rdf.DefaultBlankNodeFactory
	}

	return m
}

func (m *factoryMapper) MapBlankNode(bn rdf.BlankNode) rdf.BlankNode {
	m.mutex.Lock()

	if mapped, known := m.known[bn.Identifier]; known {
		m.mutex.Unlock()

		return rdf.BlankNode{
			Identifier: mapped,
		}
	}

	bnNext := m.factory.NewBlankNode()

	m.known[bn.Identifier] = bnNext.Identifier

	m.mutex.Unlock()

	return bnNext
}
