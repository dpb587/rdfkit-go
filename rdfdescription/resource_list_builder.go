package rdfdescription

import (
	"context"
	"iter"
	"maps"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/triples"
)

type ResourceListBuilder struct {
	resourceBySubject   map[rdf.SubjectValue]ObjectStatementList
	blankNodeReferences map[rdf.BlankNodeIdentifier]int
}

var _ triples.GraphWriter = &ResourceListBuilder{}

func NewResourceListBuilder() *ResourceListBuilder {
	return &ResourceListBuilder{
		resourceBySubject:   map[rdf.SubjectValue]ObjectStatementList{},
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
			rb.blankNodeReferences[objectSubject.Identifier]++
		}
	}
}

func (rb *ResourceListBuilder) GetBlankNodeReferences(bn rdf.BlankNode) int {
	return rb.blankNodeReferences[bn.Identifier]
}

func (rb *ResourceListBuilder) Subjects() iter.Seq[rdf.SubjectValue] {
	return maps.Keys(rb.resourceBySubject)
}

func (rb *ResourceListBuilder) GetSubjectStatements(s rdf.SubjectValue) StatementList {
	return rb.resourceBySubject[s].AsStatementList()
}

type ExportResourceOptions struct {
	// UseAnonResource will convert a SubjectResource, whose subject is a Blank Node and unreferenced by any object
	// known to ResourceListBuilder, into an AnonResource.
	UseAnonResource bool

	// Inline will convert an ObjectStatement, whose object is a Blank Node and unreferenced by any other object known
	// to ResourceListBuilder, into an AnonResourceStatement.
	//
	// Inlined resources will not be enumerated as an exported, root-level Resource.
	Inline bool
}

var DefaultExportResourceOptions = ExportResourceOptions{
	UseAnonResource: true,
	Inline:          true,
}

func (rb *ResourceListBuilder) ExportResources(opts ExportResourceOptions) iter.Seq[Resource] {
	return func(yield func(Resource) bool) {
		for subject := range rb.resourceBySubject {
			if opts.Inline {
				if bn, ok := subject.(rdf.BlankNode); ok && rb.blankNodeReferences[bn.Identifier] == 1 {
					continue
				}
			}

			if !yield(rb.ExportResource(subject, opts)) {
				return
			}
		}
	}
}

func (rb *ResourceListBuilder) ExportResource(s rdf.SubjectValue, opts ExportResourceOptions) Resource {
	statements := rb.ExportResourceStatements(s, opts)

	if opts.UseAnonResource {
		if sBlankNode, ok := s.(rdf.BlankNode); ok && rb.GetBlankNodeReferences(sBlankNode) == 0 {
			return AnonResource{
				Statements: statements,
			}
		}
	}

	return SubjectResource{
		Subject:    s,
		Statements: statements,
	}
}

func (rb *ResourceListBuilder) ExportResourceStatements(subject rdf.SubjectValue, opts ExportResourceOptions) StatementList {
	var statements StatementList

	for _, statement := range rb.resourceBySubject[subject] {
		if opts.Inline {
			if bn, ok := statement.Object.(rdf.BlankNode); ok && rb.blankNodeReferences[bn.Identifier] == 1 {
				statements = append(statements, AnonResourceStatement{
					Predicate: statement.Predicate,
					AnonResource: AnonResource{
						Statements: rb.ExportResourceStatements(bn, opts),
					},
				})

				continue
			}
		}

		statements = append(statements, statement)
	}

	return statements
}

//

func (rb *ResourceListBuilder) ToResourceWriter(ctx context.Context, e ResourceWriter, opts ExportResourceOptions) error {
	var err error

	for r := range rb.ExportResources(opts) {
		err = e.AddResource(ctx, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rb *ResourceListBuilder) ToDatasetResourceWriter(ctx context.Context, e DatasetResourceWriter, g rdf.GraphNameValue, opts ExportResourceOptions) error {
	var err error

	for r := range rb.ExportResources(opts) {
		err = e.AddDatasetResource(ctx, DatasetResource{
			Resource:  r,
			GraphName: g,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
