package pvspecification

import (
	"testing"

	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

func TestParseText(t *testing.T) {
	for _, tc := range []struct {
		Input                              string
		Expected                           ShorthandText
		ExpectedErr                        error
		ExpectedText                       string
		ExpectedPropertyValueSpecification rdfdescription.AnonResource
	}{
		// https://schema.org/docs/actions.html
		{
			Input: "required maxlength=100 name=q",
			Expected: ShorthandText{
				Attributes: []ShorthandTextAttribute{
					{
						Name: "required",
					},
					{
						Name:  "maxlength",
						Value: ptr.Value("100"),
					},
					{
						Name:  "name",
						Value: ptr.Value("q"),
					},
				},
			},
			ExpectedText: "required maxlength=100 name=q",
			ExpectedPropertyValueSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.PropertyValueSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.ValueRequired_Property,
						Object:    schemairi.True_Boolean,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.ValueMaxLength_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Text_Class,
							LexicalForm: "100",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Name_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Text_Class,
							LexicalForm: "q",
						},
					},
				},
			},
		},
		// https://schema.org/SolveMathAction#eg-0463
		{
			Input: "required name=math_expression_string",
			Expected: ShorthandText{
				Attributes: []ShorthandTextAttribute{
					{
						Name: "required",
					},
					{
						Name:  "name",
						Value: ptr.Value("math_expression_string"),
					},
				},
			},
			ExpectedText: "required name=math_expression_string",
			ExpectedPropertyValueSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.PropertyValueSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.ValueRequired_Property,
						Object:    schemairi.True_Boolean,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Name_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Text_Class,
							LexicalForm: "math_expression_string",
						},
					},
				},
			},
		},
		// validator.schema.org (2025-04-08)
		{
			Input: "required value=hello%20world",
			Expected: ShorthandText{
				Attributes: []ShorthandTextAttribute{
					{
						Name: "required",
					},
					{
						Name:  "value",
						Value: ptr.Value("hello%20world"),
					},
				},
			},
			ExpectedText: "required value=hello%20world",
			ExpectedPropertyValueSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.PropertyValueSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.ValueRequired_Property,
						Object:    schemairi.True_Boolean,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Value_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Text_Class,
							LexicalForm: "hello%20world",
						},
					},
				},
			},
		},
		// validator.schema.org (2025-04-08)
		{
			Input: "required value=\"hello world\"",
			Expected: ShorthandText{
				Attributes: []ShorthandTextAttribute{
					{
						Name: "required",
					},
					{
						Name:  "value",
						Value: ptr.Value("\"hello"),
					},
					// not shown on validator; presumably due to not matching a valid predicate
					{
						Name: "world\"",
					},
				},
			},
			ExpectedText: "required value=\"hello",
			ExpectedPropertyValueSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.PropertyValueSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.ValueRequired_Property,
						Object:    schemairi.True_Boolean,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Value_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Text_Class,
							LexicalForm: "\"hello",
						},
					},
				},
			},
		},
	} {
		t.Run(tc.Input, func(t *testing.T) {
			result := ParseText(tc.Input)
			if _a, _e := result.Err(), tc.ExpectedErr; _a != nil || _e != nil {
				if _a != nil && _e != nil {
					if _a.Error() != _e.Error() {
						t.Fatalf("expected error %v, got %v", _e, _a)
					}
				}

				t.Fatalf("expected %v, got %v", _e, _a)
			}

			resultData := result.Data

			if _a, _e := len(resultData.Attributes), len(tc.Expected.Attributes); _a != _e {
				t.Fatalf("expected %d entries, got %d", _e, _a)
			}

			for entryIdx, entry := range resultData.Attributes {
				if _a, _e := entry.Name, tc.Expected.Attributes[entryIdx].Name; _a != _e {
					t.Fatalf("entry %d: key: expected %s, got %s", entryIdx, _e, _a)
				}

				{
					_a, _e := entry.Value, tc.Expected.Attributes[entryIdx].Value
					if _a == nil && _e == nil {
						// good
					} else if _a == nil {
						t.Fatalf("entry %d: value: expected %s, got nil", entryIdx, *_e)
					} else if _e == nil {
						t.Fatalf("entry %d: value: expected nil, got %s", entryIdx, *_a)
					} else if *_a != *_e {
						t.Fatalf("entry %d: value: expected %s, got %s", entryIdx, *_e, *_a)
					}
				}
			}
		})
	}
}
