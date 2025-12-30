package triples

import "github.com/dpb587/rdfkit-go/rdf"

func Collect(iter rdf.TripleIterator) (rdf.TripleList, error) {
	var all rdf.TripleList

	for iter.Next() {
		all = append(all, iter.Triple())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return all, nil
}

func CollectErr(iter rdf.TripleIterator, err error) (rdf.TripleList, error) {
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	return Collect(iter)
}
