package encodingutil

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
)

type TextOffsetsBuilderFunc func(pairs ...any) encoding.StatementTextOffsets

func BuildTextOffsetsNil(pairs ...any) encoding.StatementTextOffsets {
	return nil
}

func BuildTextOffsetsValue(pairs ...any) encoding.StatementTextOffsets {
	v := encoding.StatementTextOffsets{}

	for i := 0; i < len(pairs); i += 2 {
		if offsets := pairs[i+1].(*cursorio.TextOffsetRange); offsets != nil {
			v[pairs[i].(encoding.StatementOffsetsType)] = *offsets
		}
	}

	return v
}
