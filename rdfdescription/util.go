package rdfdescription

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/rdfutil"
)

type builderLookup interface {
	GetBlankNodeReferences(bn rdf.BlankNode) int
}

func PreferAnonResource(b builderLookup, r Resource) Resource {
	rSubject, ok := r.(SubjectResource)
	if !ok {
		return r
	}

	bn, ok := rSubject.GetResourceSubject().(rdf.BlankNode)
	if !ok {
		return r
	}

	if b.GetBlankNodeReferences(bn) == 0 {
		return AnonResource{
			Statements: rSubject.Statements,
		}
	}

	return r
}

//

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
