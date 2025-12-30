package quads

import "github.com/dpb587/rdfkit-go/rdf"

func Collect(iter rdf.QuadIterator) (rdf.QuadList, error) {
	var all rdf.QuadList

	for iter.Next() {
		all = append(all, iter.Quad())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return all, nil
}

func CollectErr(iter rdf.QuadIterator, err error) (rdf.QuadList, error) {
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	return Collect(iter)
}
