package rdfdescription

import "github.com/dpb587/rdfkit-go/rdf/rdfutil"

func NewStatementsFromObjectsByPredicate(po rdfutil.ObjectsByPredicate) StatementList {
	var statements StatementList

	for predicate, objects := range po {
		for _, object := range objects {
			statements = append(statements, ObjectStatement{
				Predicate: predicate,
				Object:    object,
			})
		}
	}

	return statements
}
