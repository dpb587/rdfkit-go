package rdfdescription

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type Statement interface {
	NewTriples(s rdf.SubjectValue) rdf.TripleList
}

type StatementList []Statement

func (l StatementList) NewTriples(s rdf.SubjectValue) rdf.TripleList {
	var res rdf.TripleList

	for _, d := range l {
		res = append(res, d.NewTriples(s)...)
	}

	return res
}

func (l StatementList) GroupByPredicate() StatementListByPredicate {
	res := StatementListByPredicate{}

	for _, d := range l {
		switch d := d.(type) {
		case ObjectStatement:
			res[d.Predicate] = append(res[d.Predicate], d)
		case AnonResourceStatement:
			res[d.Predicate] = append(res[d.Predicate], d)
		}
	}

	return res
}

type StatementListByPredicate map[rdf.PredicateValue]StatementList

func (l StatementListByPredicate) GetPredicateList() rdf.PredicateValueList {
	var res rdf.PredicateValueList

	for k := range l {
		res = append(res, k)
	}

	return res
}

//

type ObjectStatement struct {
	Predicate rdf.PredicateValue
	Object    rdf.ObjectValue
}

var _ Statement = ObjectStatement{}

func (l ObjectStatement) NewTriples(s rdf.SubjectValue) rdf.TripleList {
	return rdf.TripleList{
		{
			Subject:   s,
			Predicate: l.Predicate,
			Object:    l.Object,
		},
	}
}

//

type AnonResourceStatement struct {
	Predicate    rdf.PredicateValue
	AnonResource AnonResource
}

var _ Statement = AnonResourceStatement{}

func (l AnonResourceStatement) NewTriples(s rdf.SubjectValue) rdf.TripleList {
	descriptionSubject, descriptionStatements := l.AnonResource.statementList()

	return append(
		descriptionStatements,
		rdf.Triple{
			Subject:   s,
			Predicate: l.Predicate,
			Object:    descriptionSubject,
		},
	)
}
