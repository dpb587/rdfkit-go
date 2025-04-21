package blanknodeutil

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/dpb587/rdfkit-go/rdf"
)

type StringerInt64Option interface {
	apply(s *builtinStringerInt64)
}

type StringerInt64Config struct {
	format    *string
	nextIndex *int64
}

// SetFormat configures the fmt-format specifier for a single, int64 argument.
//
// Default is `b%d`.
func (c StringerInt64Config) SetFormat(v string) StringerInt64Config {
	c.format = &v

	return c
}

// SetNextIndex configures the next number which will be generated. The minimum value is 0.
//
// Default is 0.
func (c StringerInt64Config) SetNextIndex(v int64) StringerInt64Config {
	if v < 0 {
		panic("value must be >= 0")
	}

	c.nextIndex = &v

	return c
}

func (c StringerInt64Config) apply(s *builtinStringerInt64) {
	if c.format != nil {
		s.format = *c.format
	}

	if c.nextIndex != nil {
		s.index.Store(*c.nextIndex - 1)
	}
}

type builtinStringerInt64 struct {
	format string
	index  *atomic.Int64

	mutex sync.Mutex
	known map[rdf.BlankNodeIdentifier]int64
}

// NewStringerInt64 generates identifiers based on an incremental int64 (minimum 0) index. It retains a reference to
// every [rdf.BlankNodeIdentifier] value it sees. It is safe for concurrent use.
func NewStringerInt64(opts ...StringerInt64Option) Stringer {
	s := &builtinStringerInt64{
		index: &atomic.Int64{},
		known: map[rdf.BlankNodeIdentifier]int64{},
	}

	s.index.Store(-1)

	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

func (s *builtinStringerInt64) GetBlankNodeIdentifier(bn rdf.BlankNode) string {
	identifier := bn.GetBlankNodeIdentifier()

	s.mutex.Lock()

	index, known := s.known[identifier]
	if !known {
		index = s.index.Add(1)

		s.known[identifier] = index
	}

	s.mutex.Unlock()

	if len(s.format) == 0 {
		return fmt.Sprintf("b%d", index)
	}

	return fmt.Sprintf(s.format, index)
}
