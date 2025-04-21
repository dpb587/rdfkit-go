package rdfioutil

import "github.com/dpb587/rdfkit-go/rdfio"

type statementIteratorIterator struct {
	iters []rdfio.StatementIterator
}

var _ rdfio.StatementIterator = &statementIteratorIterator{}

func NewStatementIteratorIterator(iters ...rdfio.StatementIterator) *statementIteratorIterator {
	return &statementIteratorIterator{
		iters: iters,
	}
}

func (i *statementIteratorIterator) Close() error {
	for _, iter := range i.iters {
		iter.Close()
	}

	return nil
}

func (i *statementIteratorIterator) Err() error {
	if len(i.iters) > 0 {
		return i.iters[0].Err()
	}

	return nil
}

func (i *statementIteratorIterator) Next() bool {
	for {
		if len(i.iters) == 0 {
			return false
		} else if i.iters[0].Err() != nil {
			return false
		} else if i.iters[0].Next() {
			return true
		}

		i.iters = i.iters[1:]
	}
}

func (i *statementIteratorIterator) GetStatement() rdfio.Statement {
	return i.iters[0].GetStatement()
}
