package openinghours

import (
	"strings"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

var dayCycle = "MoTuWeThFrSaSuMoTuWeThFrSaSu"
var dayTokens = map[string]uint8{
	"Mo": 0,
	"Tu": 2,
	"We": 4,
	"Th": 6,
	"Fr": 8,
	"Sa": 10,
	"Su": 12,
}

type PropertyValue struct {
	DayOfWeekMo bool
	DayOfWeekTu bool
	DayOfWeekWe bool
	DayOfWeekTh bool
	DayOfWeekFr bool
	DayOfWeekSa bool
	DayOfWeekSu bool

	OpensTime  string
	ClosesTime string
}

func (o PropertyValue) AsText() string {
	var tokens []string

	{
		var days = []struct {
			S string
			B bool
		}{
			{"Mo", o.DayOfWeekMo},
			{"Tu", o.DayOfWeekTu},
			{"We", o.DayOfWeekWe},
			{"Th", o.DayOfWeekTh},
			{"Fr", o.DayOfWeekFr},
			{"Sa", o.DayOfWeekSa},
			{"Su", o.DayOfWeekSu},
		}

		var tokenList []string
		var token string

		for i, day := range days {
			if day.B {
				if len(token) == 0 {
					token = day.S
				}
			} else if len(token) > 0 {
				if token == days[i-1].S {
					tokenList = append(tokenList, token)
				} else {
					tokenList = append(tokenList, token+"-"+days[i-1].S)
				}

				token = ""
			}
		}

		if len(token) > 0 {
			if token == "Su" {
				tokenList = append(tokenList, token)
			} else if o.DayOfWeekMo && (!o.DayOfWeekTu || !o.DayOfWeekWe || !o.DayOfWeekTh || !o.DayOfWeekFr || !o.DayOfWeekSa) {
				firstSplit := strings.Split(tokenList[0], "-")

				if len(firstSplit) == 2 {
					tokenList[0] = token + "-" + firstSplit[1]
				} else {
					tokenList[0] = token + "-" + "Mo"
				}
			} else {
				tokenList = append(tokenList, token+"-Su")
			}
		}

		tokens = append(tokens, strings.Join(tokenList, ","))
	}

	if len(o.OpensTime) > 0 && len(o.ClosesTime) > 0 {
		tokens = append(tokens, o.OpensTime+"-"+o.ClosesTime)
	}

	return strings.Join(tokens, " ")
}

func (o PropertyValue) AsOpeningHoursSpecification() rdfdescription.AnonResource {
	res := rdfdescription.AnonResource{
		Statements: rdfdescription.StatementList{
			rdfdescription.ObjectStatement{
				Predicate: rdfiri.Type_Property,
				Object:    schemairi.OpeningHoursSpecification_Class,
			},
		},
	}

	if o.DayOfWeekMo {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Monday_DayOfWeek,
		})
	}

	if o.DayOfWeekTu {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Tuesday_DayOfWeek,
		})
	}

	if o.DayOfWeekWe {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Wednesday_DayOfWeek,
		})
	}

	if o.DayOfWeekTh {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Thursday_DayOfWeek,
		})
	}

	if o.DayOfWeekFr {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Friday_DayOfWeek,
		})
	}

	if o.DayOfWeekSa {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Saturday_DayOfWeek,
		})
	}

	if o.DayOfWeekSu {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.DayOfWeek_Property,
			Object:    schemairi.Sunday_DayOfWeek,
		})
	}

	if len(o.OpensTime) > 0 {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.Opens_Property,
			Object: rdf.Literal{
				Datatype:    schemairi.Time_Class,
				LexicalForm: o.OpensTime,
			},
		})
	}

	if len(o.ClosesTime) > 0 {
		res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
			Predicate: schemairi.Closes_Property,
			Object: rdf.Literal{
				Datatype:    schemairi.Time_Class,
				LexicalForm: o.ClosesTime,
			},
		})
	}

	return res
}
