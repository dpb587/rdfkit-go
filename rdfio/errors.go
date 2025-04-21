package rdfio

import "errors"

var (
	ErrNodeNotBound      = errors.New("node not bound")
	ErrStatementNotBound = errors.New("statement not bound")
)

type StatementError struct {
	Statement Statement
	Err       error
}

func (e StatementError) Error() string {
	return e.Err.Error()
}
