package blanknodeutil

import (
	"sync"
	"sync/atomic"

	"github.com/dpb587/rdfkit-go/rdf"
)

// StringMapper maps a string identifier to a blank node. For a given identifier, the same blank node will always be
// returned.
//
// The `_:` prefix used by many encodings to begin a blank node should not be included in the string. That is, in a
// Turtle document, a `_:b0` term should be passed as `b0`.
type StringMapper interface {
	MapBlankNodeIdentifier(v string) rdf.BlankNode
}

//

type builtinStringTable struct {
	factory *builtinFactory
	mutex   sync.Mutex
	known   map[string]int64
}

var _ StringMapper = (*builtinStringTable)(nil)

func NewStringMapper() StringMapper {
	st := &builtinStringTable{
		factory: &builtinFactory{
			a: &atomic.Int64{},
		},
		known: map[string]int64{},
	}

	return st
}

func (st *builtinStringTable) MapBlankNodeIdentifier(v string) rdf.BlankNode {
	if len(v) == 0 {
		return rdf.NewBlankNodeWithIdentifier(
			builtinFactoryIdentifier{
				g: st.factory,
				i: st.factory.a.Add(1),
			},
		)
	}

	st.mutex.Lock()

	i, known := st.known[v]
	if !known {
		i = st.factory.a.Add(1)

		st.known[v] = i
	}

	st.mutex.Unlock()

	return rdf.NewBlankNodeWithIdentifier(
		builtinFactoryIdentifier{
			g: st.factory,
			i: i,
		},
	)
}
