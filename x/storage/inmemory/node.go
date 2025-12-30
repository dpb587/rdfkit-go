package inmemory

import (
	"github.com/dpb587/rdfkit-go/rdf"
)

type Node struct {
	d *Dataset
	t rdf.Term

	Baggage map[any]any
}

func (br *Node) GetDataset() *Dataset {
	return br.d
}

func (br *Node) GetTerm() rdf.Term {
	return br.t
}

//

type nodeList []*Node
