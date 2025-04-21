package rdfio

type StatementIterator interface {
	Close() error
	Err() error
	Next() bool
	GetStatement() Statement
}

func CollectStatementsErr(iter StatementIterator, err error) (StatementList, error) {
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	return CollectStatements(iter)
}

func CollectStatements(iter StatementIterator) (StatementList, error) {
	var all StatementList

	for iter.Next() {
		all = append(all, iter.GetStatement())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return all, nil
}
