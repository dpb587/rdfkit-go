package literalutil

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type CustomValue interface {
	termutil.CustomValue

	AsLiteralTerm() rdf.Literal
}
