package encoding

import (
	"github.com/dpb587/cursorio-go/cursorio"
)

type StatementOffsetsType int

const (
	GraphNameStatementOffsets StatementOffsetsType = iota
	SubjectStatementOffsets
	PredicateStatementOffsets
	ObjectStatementOffsets
)

func StatementOffsetsTypeName(t StatementOffsetsType) string {
	switch t {
	case GraphNameStatementOffsets:
		return "graphName"
	case SubjectStatementOffsets:
		return "subject"
	case PredicateStatementOffsets:
		return "predicate"
	case ObjectStatementOffsets:
		return "object"
	}

	return "unknown"
}

//

type StatementTextOffsets map[StatementOffsetsType]cursorio.TextOffsetRange

//

type TextOffsetsStatement interface {
	GetStatementTextOffsets() StatementTextOffsets
}
