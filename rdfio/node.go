package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type Node interface {
	GetTerm() rdf.Term
}

type NodeList []Node
