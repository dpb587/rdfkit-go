package encodingtest

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

func CollectTripleStatements(iter rdf.TripleIterator) (TripleStatementList, error) {
	var statements TripleStatementList

	iterTextStatements, _ := iter.(encoding.StatementTextOffsetsProvider)

	for iter.Next() {
		quadStatement := TripleStatement{
			Triple: iter.Triple(),
		}

		if iterTextStatements != nil {
			quadStatement.TextOffsets = iterTextStatements.StatementTextOffsets()
		}

		statements = append(statements, quadStatement)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

func CollectTripleStatementsErr(iter rdf.TripleIterator, err error) (TripleStatementList, error) {
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	return CollectTripleStatements(iter)
}
