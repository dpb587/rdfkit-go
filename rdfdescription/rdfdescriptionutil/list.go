package rdfdescriptionutil

import (
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

func NewObjectValueListStatement(predicate rdf.PredicateValue, values ...rdf.ObjectValue) rdfdescription.Statement {
	vl := len(values)

	if vl == 0 {
		return rdfdescription.ObjectStatement{
			Predicate: predicate,
			Object:    rdfiri.Nil_List,
		}
	}

	vl -= 1

	listResource := rdfdescription.AnonResource{
		Statements: rdfdescription.StatementList{
			rdfdescription.ObjectStatement{
				Predicate: rdfiri.First_Property,
				Object:    values[vl],
			},
			rdfdescription.ObjectStatement{
				Predicate: rdfiri.Rest_Property,
				Object:    rdfiri.Nil_List,
			},
		},
	}

	for vl -= 1; vl >= 0; vl-- {
		listResource = rdfdescription.AnonResource{
			Statements: rdfdescription.StatementList{
				rdfdescription.ObjectStatement{
					Predicate: rdfiri.First_Property,
					Object:    values[vl],
				},
				rdfdescription.AnonResourceStatement{
					Predicate:    rdfiri.Rest_Property,
					AnonResource: listResource,
				},
			},
		}
	}

	return rdfdescription.AnonResourceStatement{
		Predicate:    predicate,
		AnonResource: listResource,
	}
}
