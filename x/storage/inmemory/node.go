package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type Node struct {
	d *Dataset
	t rdf.Term

	Baggage map[any]any
}

var _ rdfio.Node = &Node{}

func (br *Node) GetDataset() rdfio.Dataset {
	return br.d
}

func (br *Node) GetTerm() rdf.Term {
	return br.t
}

//

type nodeList []*Node
