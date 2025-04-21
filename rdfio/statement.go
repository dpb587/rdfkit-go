package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type Statement interface {
	GetGraphName() rdf.GraphNameValue
	GetTriple() rdf.Triple
}

type StatementList []Statement
