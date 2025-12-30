package xsdtype

import (
	"fmt"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

type Short int16

var _ objecttypes.Value = Short(0)

func MapShort(lexicalForm string) (Short, error) {
	vInt64, err := strconv.ParseInt(xsdutil.WhiteSpaceCollapse(lexicalForm), 10, 16)
	if err != nil {
		return Short(0), fmt.Errorf("%w: %v", rdf.ErrLiteralLexicalFormNotValid, err)
	}

	return Short(vInt64), nil
}

func (v Short) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    xsdiri.Short_Datatype,
		LexicalForm: strconv.FormatInt(int64(v), 10),
	}
}

func (Short) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Short) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Short_Datatype {
		return false
	}

	return strconv.FormatInt(int64(v), 10) == tLiteral.LexicalForm
}
