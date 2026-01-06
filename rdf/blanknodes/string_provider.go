package blanknodes

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

// StringProvider provides a string for a blank node. For a given identifier, the same string identifier must always be
// returned.
//
// The `_:` prefix should not be included in the string.
type StringProvider interface {
	GetBlankNodeString(bn rdf.BlankNode) string
}

// StringProviderProvider may be implemented by types that can provide a StringProvider.
//
// Specifically, it may be offered by StringFactory implementations if they can convert their BlankNode values back to a
// string.
//
// If a BlankNode is considered out of range, fallback may be used instead.
type StringProviderProvider interface {
	GetStringProvider(fallback StringProvider) StringProvider
}
