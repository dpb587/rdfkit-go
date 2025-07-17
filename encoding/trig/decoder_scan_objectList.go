package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func reader_scan_ObjectList_Continue(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_objectList.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	}

	if r0.Rune == ',' {
		r.commit(r0.AsDecodedRunes())

		r.pushState(ectx, reader_scan_ObjectList_Continue)

		return readerStack{ectx, reader_scan_Object}, nil
	}

	r.buf.BacktrackRunes(r0)

	return readerStack{}, nil
}
