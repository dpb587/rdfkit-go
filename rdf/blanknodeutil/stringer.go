package blanknodeutil

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

// Stringer provides a string identifier for a blank node. For a given blank node, the same string identifier will
// always be returned.
//
// The `_:` prefix used by many encodings to begin a blank node should not be included in the string.
type Stringer interface {
	GetBlankNodeIdentifier(t rdf.BlankNode) string
}
