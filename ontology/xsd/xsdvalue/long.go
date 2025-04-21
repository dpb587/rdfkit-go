package xsdvalue

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

type Long int64

var _ literalutil.CustomValue = Long(0)

func MapLong(lexicalForm string) (Long, error) {
	vInt64, err := strconv.ParseInt(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 16)
	if err != nil {
		return Long(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Long(vInt64), nil
}

func (v Long) AsLiteralTerm() rdf.Literal {
	return rdf.Literal{
		Datatype:    xsdiri.Long_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (Long) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Long) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Long_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
