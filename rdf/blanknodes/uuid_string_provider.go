package blanknodes

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/google/uuid"
)

type uuidStringProvider struct {
	format string
	reader io.Reader
	mutex  sync.Mutex
	known  map[rdf.BlankNodeIdentifier]uuid.UUID
}

var _ StringProvider = &uuidStringProvider{}

// NewUUIDStringProvider uniquely maps a BlankNode to a UUID-based string.
//
// The default format is "%s", and reader is crypto/rand.Reader.
//
// If reader errors, the code panics.
//
// This retains a reference to all BlankNodeIdentifier values.
func NewUUIDStringProvider(format string, reader io.Reader) StringProvider {
	if len(format) == 0 {
		format = "%s"
	}

	return &uuidStringProvider{
		format: format,
		reader: rand.Reader,
		known:  make(map[rdf.BlankNodeIdentifier]uuid.UUID),
	}
}

func (sp *uuidStringProvider) GetBlankNodeString(bn rdf.BlankNode) string {
	sp.mutex.Lock()

	index, known := sp.known[bn.Identifier]
	if !known {
		value, err := uuid.NewV7FromReader(sp.reader)
		if err != nil {
			panic(fmt.Errorf("%T: %v", sp, err))
		}

		sp.known[bn.Identifier] = value
	}

	sp.mutex.Unlock()

	return fmt.Sprintf(sp.format, index)
}
