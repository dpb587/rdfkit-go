package terms

import "github.com/dpb587/rdfkit-go/rdf"

type Formatter interface {
	FormatTerm(t rdf.Term) string
}
