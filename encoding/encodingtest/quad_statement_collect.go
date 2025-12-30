package encodingtest

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

func CollectQuadStatements(iter rdf.QuadIterator) (QuadStatementList, error) {
	var statements QuadStatementList

	iterTextStatements, _ := iter.(encoding.StatementTextOffsetsProvider)

	for iter.Next() {
		quadStatement := QuadStatement{
			Quad: iter.Quad(),
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

func CollectQuadStatementsErr(iter rdf.QuadIterator, err error) (QuadStatementList, error) {
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	return CollectQuadStatements(iter)
}
