package rdfdescription

import "github.com/dpb587/rdfkit-go/rdf"

type DatasetResource struct {
	GraphName rdf.GraphNameValue
	Resource  Resource
}

func (dr DatasetResource) NewQuads() rdf.QuadList {
	return dr.Resource.NewTriples().AsQuads(dr.GraphName)
}
