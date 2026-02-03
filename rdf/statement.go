package rdf

type StatementType int

const (
	TripleStatementType StatementType = iota
	QuadStatementType
)

//

// Statement represents the tuple of an RDF statement.
//
// This is a closed interface. See [Quad] and [Triple].
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
