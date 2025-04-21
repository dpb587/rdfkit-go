package rdfio

type NodeIterator interface {
	Close() error
	Err() error
	Next() bool
	GetNode() Node
}

func EnumerateNodes(iter NodeIterator) (NodeList, error) {
	var all NodeList

	for iter.Next() {
		all = append(all, iter.GetNode())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return all, nil
}
