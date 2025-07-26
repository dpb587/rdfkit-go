package schemairi

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

var AlternateBaseList = []rdf.IRI{
	"https://schema.org/",
	"http://www.schema.org/",
	"https://www.schema.org/",
}

func NormalizeBase(b rdf.IRI) rdf.IRI {
	for _, alt := range AlternateBaseList {
		if len(b) >= len(alt) && b[:len(alt)] == alt {
			return Base + b[len(alt):]
		}
	}

	return b
}
