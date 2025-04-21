package schemaliteral

import (
	"fmt"
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdliteral"
	"github.com/dpb587/rdfkit-go/rdf"
)

func TestCastTime(t *testing.T) {
	for _, tc := range []struct {
		Input               rdf.ObjectValue
		ExpectedObjectValue rdf.ObjectValue
		ExpectedOK          bool
	}{
		{
			Input: xsdliteral.NewString("09:00"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Time_Class,
				LexicalForm: "09:00",
			},
			ExpectedOK: true,
		},
		{
			Input: xsdliteral.NewString("09:00:00"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Time_Class,
				LexicalForm: "09:00:00",
			},
			ExpectedOK: true,
		},
		{
			Input: xsdliteral.NewString("9:00 AM"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Time_Class,
				LexicalForm: "09:00",
			},
			ExpectedOK: true,
		},
		{
			Input: xsdliteral.NewString(" 9:00 AM"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Time_Class,
				LexicalForm: "09:00",
			},
			ExpectedOK: true,
		},
	} {
		t.Run(fmt.Sprintf("%v", tc.Input), func(t *testing.T) {
			objectValue, ok := CastTime(tc.Input, CastOptions{})
			if _a, _e := ok, tc.ExpectedOK; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := objectValue, tc.ExpectedObjectValue; !_a.TermEquals(_e) {
				t.Fatalf("expected %v, got %v", _e, _a)
			}
		})
	}
}
