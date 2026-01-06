package blanknodes

import "github.com/dpb587/rdfkit-go/rdf"

// StringProvider provides a string for a blank node. For a given identifier, the same string identifier must always be
// returned.
//
// The `_:` prefix should not be included in the string.
type StringProvider interface {
	GetBlankNodeString(bn rdf.BlankNode) string
}
