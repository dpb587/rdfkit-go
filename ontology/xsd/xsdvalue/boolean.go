package xsdvalue

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/literalutil"
)

var (
	booleanTrue = rdf.Literal{
		Datatype:    xsdiri.Boolean_Datatype,
		LexicalForm: "true",
	}
	booleanTrueValue = Boolean(true)

	booleanFalse = rdf.Literal{
		Datatype:    xsdiri.Boolean_Datatype,
		LexicalForm: "false",
	}
	booleanFalseValue = Boolean(false)
)

type Boolean bool

var _ literalutil.CustomValue = Boolean(false)

func MapBoolean(lexicalForm string) (Boolean, error) {
	switch xsdutil.WhiteSpaceCollapse(lexicalForm) {
	case "true", "1":
		return booleanTrueValue, nil
	case "false", "0":
		return booleanFalseValue, nil
	}

	return Boolean(false), rdf.ErrLiteralLexicalFormNotValid
}

func (v Boolean) AsLiteralTerm() rdf.Literal {
	if v {
		return booleanTrue
	}

	return booleanFalse
}

func (Boolean) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Boolean) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Boolean_Datatype {
		return false
	} else if bool(v) {
		return tLiteral.LexicalForm == "true"
	}

	return tLiteral.LexicalForm == "false"
}
