package rdfioutil

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfio"
)

func CompareStatementsDeterministic(i, j rdfio.Statement) int {
	iop, ok := i.(encoding.DecoderTextOffsetsStatement)
	if !ok {
		return 0
	}

	jop, ok := j.(encoding.DecoderTextOffsetsStatement)
	if !ok {
		return 0
	}

	io := iop.GetDecoderTextOffsets()
	jo := jop.GetDecoderTextOffsets()

	if io == nil && jo == nil {
		return 0
	} else if io == nil {
		return -1
	} else if jo == nil {
		return 1
	}

	for _, t := range []encoding.StatementOffsetsType{
		encoding.ObjectStatementOffsets,
		encoding.PredicateStatementOffsets,
		encoding.SubjectStatementOffsets,
	} {
		ib, ibok := io[t]
		jb, jbok := jo[t]

		if !ibok && !jbok {
			continue
		} else if !ibok {
			return -1
		} else if !jbok {
			return 1
		}

		if ib.From.Byte < jb.From.Byte {
			return -1
		} else if ib.From.Byte > jb.From.Byte {
			return 1
		}
	}

	return 0
}
