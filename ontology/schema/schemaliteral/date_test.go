package schemaliteral

import (
	"fmt"
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdliteral"
	"github.com/dpb587/rdfkit-go/rdf"
)

func TestCastDate(t *testing.T) {
	for _, tc := range []struct {
		Input               rdf.ObjectValue
		ExpectedObjectValue rdf.ObjectValue
		ExpectedOK          bool
	}{
		{
			Input: xsdliteral.NewString("2020-04-01"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Date_Class,
				LexicalForm: "2020-04-01",
			},
			ExpectedOK: true,
		},
		{
			Input: xsdliteral.NewString("2020-04-01-07:00"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Date_Class,
				LexicalForm: "2020-04-01-07:00",
			},
			ExpectedOK: true,
		},
		{
			// expected normalization of 00:00 to Z
			Input: xsdliteral.NewString("2020-04-01-00:00"),
			ExpectedObjectValue: rdf.Literal{
				Datatype:    schemairi.Date_Class,
				LexicalForm: "2020-04-01Z",
			},
			ExpectedOK: true,
		},
	} {
		t.Run(fmt.Sprintf("%v", tc.Input), func(t *testing.T) {
			objectValue, ok := CastDate(tc.Input, CastOptions{})
			if _a, _e := ok, tc.ExpectedOK; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := objectValue, tc.ExpectedObjectValue; !_a.TermEquals(_e) {
				t.Fatalf("expected %v, got %v", _e, _a)
			}
		})
	}
}
