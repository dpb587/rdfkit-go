package rdfdescription

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/triples"
)

type ResourceListBuilder struct {
	resourceBySubject   map[rdf.SubjectValue][]ObjectStatement
	blankNodeReferences map[rdf.BlankNodeIdentifier]int
}

var _ triples.GraphWriter = &ResourceListBuilder{}

func NewResourceListBuilder() *ResourceListBuilder {
	return &ResourceListBuilder{
		resourceBySubject:   map[rdf.SubjectValue][]ObjectStatement{},
		blankNodeReferences: map[rdf.BlankNodeIdentifier]int{},
	}
}

// AddTriple is a wrapper for Add to satisfy the triples.StorageWriter interface.
func (rb *ResourceListBuilder) AddTriple(_ context.Context, t rdf.Triple) error {
	rb.Add(t)

	return nil
}

func (rb *ResourceListBuilder) Add(triples ...rdf.Triple) {
	for _, t := range triples {
		rb.resourceBySubject[t.Subject] = append(rb.resourceBySubject[t.Subject], ObjectStatement{
			Predicate: t.Predicate,
			Object:    t.Object,
		})

		switch objectSubject := t.Object.(type) {
		case rdf.BlankNode:
			rb.blankNodeReferences[objectSubject.GetBlankNodeIdentifier()]++
		}
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

func (rb *ResourceListBuilder) AddTo(ctx context.Context, e ResourceWriter, preferAnon bool) error {
	var err error

	for _, r := range rb.GetResources() {
		if preferAnon {
			err = e.AddResource(ctx, PreferAnonResource(rb, r))
		} else {
			err = e.AddResource(ctx, r)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (rb *ResourceListBuilder) AddToDataset(ctx context.Context, e DatasetResourceWriter, g rdf.GraphNameValue, preferAnon bool) error {
	var err error

	for _, r := range rb.GetResources() {
		if preferAnon {
			err = e.AddDatasetResource(ctx, PreferAnonResource(rb, r), g)
		} else {
			err = e.AddDatasetResource(ctx, r, g)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
