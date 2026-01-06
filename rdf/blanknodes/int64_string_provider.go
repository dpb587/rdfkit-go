package blanknodes

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/dpb587/rdfkit-go/rdf"
)

type int64StringProvider struct {
	format string
	value  *atomic.Int64
	mutex  sync.Mutex
	known  map[rdf.BlankNodeIdentifier]int64
}

var _ StringProvider = &int64StringProvider{}

// NewInt64StringProvider uniquely maps a BlankNode to an int64-based string. An empty format defaults to "b%d".
//
// This retains a reference to all BlankNodeIdentifier values.
func NewInt64StringProvider(format string) StringProvider {
	value := &atomic.Int64{}
	value.Store(-1)

	if len(format) == 0 {
		format = "b%d"
	}

	return &int64StringProvider{
		format: format,
		value:  value,
		known:  make(map[rdf.BlankNodeIdentifier]int64),
	}
}

func (sp *int64StringProvider) GetBlankNodeString(bn rdf.BlankNode) string {
	sp.mutex.Lock()

	index, known := sp.known[bn.Identifier]
	if !known {
		index = sp.value.Add(1)

		sp.known[bn.Identifier] = index
	}

	sp.mutex.Unlock()

	return fmt.Sprintf(sp.format, index)
}
