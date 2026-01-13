package rdfdescription

import (
	"context"
	"maps"
	"slices"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/quads"
)

type DatasetResource struct {
	GraphName rdf.GraphNameValue
	Resource  Resource
}

func (dr DatasetResource) NewQuads() rdf.QuadList {
	return dr.Resource.NewTriples().AsQuads(dr.GraphName)
}

//

type DatasetResourceListBuilder struct {
	builderByGraphName map[rdf.GraphNameValue]*ResourceListBuilder
}

var _ quads.DatasetWriter = &DatasetResourceListBuilder{}

func NewDatasetResourceListBuilder() *DatasetResourceListBuilder {
	return &DatasetResourceListBuilder{
		builderByGraphName: map[rdf.GraphNameValue]*ResourceListBuilder{},
	}
}

func (e *DatasetResourceListBuilder) GetGraphNames() rdf.GraphNameValueList {
	return slices.Collect(maps.Keys(e.builderByGraphName))
}

func (e *DatasetResourceListBuilder) GetResourceListBuilder(graphName rdf.GraphNameValue) *ResourceListBuilder {
	return e.builderByGraphName[graphName]
}

func (rb *DatasetResourceListBuilder) GetBlankNodeReferences(bn rdf.BlankNode) int {
	var count int

	for _, builder := range rb.builderByGraphName {
		count += builder.GetBlankNodeReferences(bn)
	}

	return count
}

func (e *DatasetResourceListBuilder) AddQuad(ctx context.Context, quad rdf.Quad) error {
	e.Add(quad)

	return nil
}

func (e *DatasetResourceListBuilder) Add(quads ...rdf.Quad) {
	for _, quad := range quads {
		if e.builderByGraphName[quad.GraphName] == nil {
			e.builderByGraphName[quad.GraphName] = NewResourceListBuilder()
		}

		e.builderByGraphName[quad.GraphName].Add(quad.Triple)
	}
}

func (e *DatasetResourceListBuilder) AddToDataset(ctx context.Context, z DatasetResourceWriter, opts ExportResourceOptions) error {
	for graphName, builder := range e.builderByGraphName {
		err := builder.AddToDataset(ctx, z, graphName, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
