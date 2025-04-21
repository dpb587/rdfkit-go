package blanknodeutil

import (
	"sync"

	"github.com/dpb587/rdfkit-go/rdf"
)

// Mapper maps blank nodes from one implementation to another.
type Mapper interface {
	MapBlankNode(t rdf.BlankNode) rdf.BlankNode
}

//

type builtinMapper struct {
	factory rdf.BlankNodeFactory
	mutex   sync.Mutex
	known   map[rdf.BlankNodeIdentifier]rdf.BlankNode
}

func NewMapper(factory rdf.BlankNodeFactory) Mapper {
	m := &builtinMapper{
		factory: factory,
		known:   map[rdf.BlankNodeIdentifier]rdf.BlankNode{},
	}

	if m.factory == nil {
		m.factory = rdf.DefaultBlankNodeFactory
	}

	return m
}

func (m *builtinMapper) MapBlankNode(t rdf.BlankNode) rdf.BlankNode {
	identifier := t.GetBlankNodeIdentifier()

	m.mutex.Lock()

	if mapped, known := m.known[identifier]; known {
		m.mutex.Unlock()

		return mapped
	}

	bn := m.factory.NewBlankNode()

	m.known[identifier] = bn

	m.mutex.Unlock()

	return bn
}
