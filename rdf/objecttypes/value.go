package objecttypes

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type Value interface {
	TermKind() rdf.TermKind
	TermEquals(a rdf.Term) bool

	AsObjectValue() rdf.ObjectValue
}
