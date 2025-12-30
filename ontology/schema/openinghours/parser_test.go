package openinghours

import (
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingassert"
)

func TestParseText(t *testing.T) {
	for _, tc := range []struct {
		Input                             string
		Expected                          PropertyValue
		ExpectedErr                       error
		ExpectedText                      string
		ExpectedOpeningHoursSpecification rdfdescription.AnonResource
	}{
		// https://schema.org/openingHours
		{
			Input: "Tu,Th 16:00-20:00",
			Expected: PropertyValue{
				DayOfWeekTu: true,
				DayOfWeekTh: true,
				OpensTime:   "16:00",
				ClosesTime:  "20:00",
			},
			ExpectedText: "Tu,Th 16:00-20:00",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Tuesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Thursday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Opens_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "16:00",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Closes_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "20:00",
						},
					},
				},
			},
		},
		{
			Input: "Mo-Su",
			Expected: PropertyValue{
				DayOfWeekMo: true,
				DayOfWeekTu: true,
				DayOfWeekWe: true,
				DayOfWeekTh: true,
				DayOfWeekFr: true,
				DayOfWeekSa: true,
				DayOfWeekSu: true,
			},
			ExpectedText: "Mo-Su",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Monday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Tuesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Wednesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Thursday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Friday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Saturday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Sunday_DayOfWeek,
					},
				},
			},
		},
		// https://schema.org/openingHours#eg-0194
		{
			Input: "Mo,Tu,We,Th 09:00-12:00",
			Expected: PropertyValue{
				DayOfWeekMo: true,
				DayOfWeekTu: true,
				DayOfWeekWe: true,
				DayOfWeekTh: true,
				OpensTime:   "09:00",
				ClosesTime:  "12:00",
			},
			ExpectedText: "Mo-Th 09:00-12:00",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Monday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Tuesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Wednesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Thursday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Opens_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "09:00",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Closes_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "12:00",
						},
					},
				},
			},
		},
		// https://schema.org/openingHours#eg-0432 (1)
		{
			Input: "Mo-Fr 10:00-19:00",
			Expected: PropertyValue{
				DayOfWeekMo: true,
				DayOfWeekTu: true,
				DayOfWeekWe: true,
				DayOfWeekTh: true,
				DayOfWeekFr: true,
				OpensTime:   "10:00",
				ClosesTime:  "19:00",
			},
			ExpectedText: "Mo-Fr 10:00-19:00",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Monday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Tuesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Wednesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Thursday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Friday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Opens_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "10:00",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Closes_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "19:00",
						},
					},
				},
			},
		},
		// https://schema.org/openingHours#eg-0432 (2)
		{
			Input: "Sa 10:00-22:00",
			Expected: PropertyValue{
				DayOfWeekSa: true,
				OpensTime:   "10:00",
				ClosesTime:  "22:00",
			},
			ExpectedText: "Sa 10:00-22:00",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Saturday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Opens_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "10:00",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Closes_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "22:00",
						},
					},
				},
			},
		},
		// https://schema.org/openingHours#eg-0432 (3)
		{
			Input: "Su 10:00-21:00",
			Expected: PropertyValue{
				DayOfWeekSu: true,
				OpensTime:   "10:00",
				ClosesTime:  "21:00",
			},
			ExpectedText: "Su 10:00-21:00",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Sunday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Opens_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "10:00",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Closes_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "21:00",
						},
					},
				},
			},
		},
		// logic: concise wraps across week end/start
		{
			Input: "We-Mo 10:00-22:00",
			Expected: PropertyValue{
				DayOfWeekMo: true,
				DayOfWeekWe: true,
				DayOfWeekTh: true,
				DayOfWeekFr: true,
				DayOfWeekSa: true,
				DayOfWeekSu: true,
				OpensTime:   "10:00",
				ClosesTime:  "22:00",
			},
			ExpectedText: "We-Mo 10:00-22:00",
			ExpectedOpeningHoursSpecification: rdfdescription.AnonResource{
				Statements: rdfdescription.StatementList{
					rdfdescription.ObjectStatement{
						Predicate: rdfiri.Type_Property,
						Object:    schemairi.OpeningHoursSpecification_Class,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Monday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Wednesday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Thursday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Friday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Saturday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.DayOfWeek_Property,
						Object:    schemairi.Sunday_DayOfWeek,
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Opens_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "10:00",
						},
					},
					rdfdescription.ObjectStatement{
						Predicate: schemairi.Closes_Property,
						Object: rdf.Literal{
							Datatype:    schemairi.Time_Class,
							LexicalForm: "22:00",
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

			if _a, _e := resultData.DayOfWeekMo, tc.Expected.DayOfWeekMo; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.DayOfWeekTu, tc.Expected.DayOfWeekTu; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.DayOfWeekWe, tc.Expected.DayOfWeekWe; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.DayOfWeekTh, tc.Expected.DayOfWeekTh; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.DayOfWeekFr, tc.Expected.DayOfWeekFr; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.DayOfWeekSa, tc.Expected.DayOfWeekSa; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.DayOfWeekSu, tc.Expected.DayOfWeekSu; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.OpensTime, tc.Expected.OpensTime; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.ClosesTime, tc.Expected.ClosesTime; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			} else if _a, _e := resultData.AsText(), tc.ExpectedText; _a != _e {
				t.Fatalf("expected %v, got %v", _e, _a)
			}

			testingassert.IsomorphicGraphs(t, tc.ExpectedOpeningHoursSpecification.NewTriples(), resultData.AsOpeningHoursSpecification().NewTriples())
		})
	}
}
