package rdfdescription

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

// TODO rdfdescriptionutil?

type ResourceListBuilder struct {
	resourceBySubject   map[rdf.SubjectValue][]ObjectStatement
	blankNodeReferences map[rdf.BlankNodeIdentifier]int
}

func NewResourceListBuilder() *ResourceListBuilder {
	return &ResourceListBuilder{
		resourceBySubject:   map[rdf.SubjectValue][]ObjectStatement{},
		blankNodeReferences: map[rdf.BlankNodeIdentifier]int{},
	}
}

func (rb *ResourceListBuilder) AddTriple(t rdf.Triple) {
	var rs = ObjectStatement{
		Predicate: t.Predicate,
		Object:    t.Object,
	}

	rb.resourceBySubject[t.Subject] = append(rb.resourceBySubject[t.Subject], rs)

	switch objectSubject := t.Object.(type) {
	case rdf.BlankNode:
		rb.blankNodeReferences[objectSubject.GetBlankNodeIdentifier()]++
	}
}

func (rb *ResourceListBuilder) GetBlankNodeReferences(bn rdf.BlankNode) int {
	return rb.blankNodeReferences[bn.GetBlankNodeIdentifier()]
}

func (rb *ResourceListBuilder) GetResourceStatements(s rdf.SubjectValue) StatementList {
	return rb.getResourceStatements(s)
}

func (rb *ResourceListBuilder) GetResources() ResourceList {
	var resources ResourceList

	for subject := range rb.resourceBySubject {
		if bn, ok := subject.(rdf.BlankNode); ok && rb.blankNodeReferences[bn.GetBlankNodeIdentifier()] == 1 {
			continue
		}

		resources = append(resources, SubjectResource{
			Subject:    subject,
			Statements: rb.getResourceStatements(subject),
		})
	}

	return resources
}

func (rb *ResourceListBuilder) GetResource(s rdf.SubjectValue) (Resource, bool) {
	if _, ok := rb.resourceBySubject[s]; !ok {
		return nil, false
	}

	return SubjectResource{
		Subject:    s,
		Statements: rb.getResourceStatements(s),
	}, true
}

func (rb *ResourceListBuilder) getResourceStatements(subject rdf.SubjectValue) StatementList {
	var statements StatementList

	for _, statement := range rb.resourceBySubject[subject] {
		if bn, ok := statement.Object.(rdf.BlankNode); ok && rb.blankNodeReferences[bn.GetBlankNodeIdentifier()] == 1 {
			statements = append(statements, AnonResourceStatement{
				Predicate: statement.Predicate,
				AnonResource: AnonResource{
					Statements: rb.getResourceStatements(bn),
				},
			})
		} else {
			statements = append(statements, statement)
		}
	}

	return statements
}
