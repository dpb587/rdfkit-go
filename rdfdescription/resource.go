package rdfdescription

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type Resource interface {
	// GetResourceSubject returns the subject of the resource. An anonymous resource will return nil.
	GetResourceSubject() rdf.SubjectValue

	// GetResourceStatements returns the statements associated with the resource.
	GetResourceStatements() StatementList

	NewTriples() rdf.TripleList
}

type ResourceList []Resource

func (rl ResourceList) NewTriples() rdf.TripleList {
	var triples rdf.TripleList

	for _, r := range rl {
		triples = append(triples, r.NewTriples()...)
	}

	return triples
}

//

type SubjectResource struct {
	// Subject is the resource term. If nil, this acts like [AnonResource].
	Subject    rdf.SubjectValue
	Statements StatementList
}

var _ Resource = (*SubjectResource)(nil)

func (d SubjectResource) GetResourceSubject() rdf.SubjectValue {
	return d.Subject
}

func (d SubjectResource) GetResourceStatements() StatementList {
	return d.Statements
}

func (d SubjectResource) NewTriples() rdf.TripleList {
	_, tb := d.statementList()

	return tb
}

func (d SubjectResource) statementList() (rdf.SubjectValue, rdf.TripleList) {
	var s = d.Subject

	if s == nil {
		s = rdf.NewBlankNode()
	}

	return s, d.Statements.NewTriples(s)
}

//

type AnonResource struct {
	Statements StatementList
}

var _ Resource = (*AnonResource)(nil)

func (d AnonResource) GetResourceSubject() rdf.SubjectValue {
	return nil
}

func (d AnonResource) GetResourceStatements() StatementList {
	return d.Statements
}

func (r AnonResource) NewTriples() rdf.TripleList {
	_, tb := r.statementList()

	return tb
}

func (r AnonResource) statementList() (rdf.SubjectValue, rdf.TripleList) {
	s := rdf.NewBlankNode()

	return s, r.Statements.NewTriples(s)
}
