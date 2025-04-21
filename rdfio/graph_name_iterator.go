package rdfio

import "github.com/dpb587/rdfkit-go/rdf"

type GraphNameIterator interface {
	Close() error
	Err() error
	Next() bool
	GetGraphName() rdf.GraphNameValue
}

func EnumerateGraphNames(iter GraphNameIterator) (rdf.GraphNameValueList, error) {
	var all rdf.GraphNameValueList

	for iter.Next() {
		all = append(all, iter.GetGraphName())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return all, nil
}
