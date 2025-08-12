package schemavalue

import (
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type Number float64

var _ termutil.CustomValue = Number(0)

func (v Number) AsTerm() rdf.Term {
	return rdf.Literal{
		Datatype:    schemairi.Number_Class,
		LexicalForm: strconv.FormatFloat(float64(v), 'f', -1, 64),
	}
}

func (Number) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Number) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != schemairi.Number_Class {
		return false
	}

	return strconv.FormatFloat(float64(v), 'f', -1, 64) == tLiteral.LexicalForm
}
