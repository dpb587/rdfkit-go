package rdf

type StatementType int

const (
	TripleStatementType StatementType = iota
	QuadStatementType
)

//

type Statement interface {
	StatementType() StatementType

	isStatement()
}

//

type StatementIterator interface {
	Next() bool
	Err() error
	Close() error

	Statement() Statement
}
