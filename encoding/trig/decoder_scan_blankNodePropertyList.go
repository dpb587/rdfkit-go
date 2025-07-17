package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func reader_scan_blankNodePropertyList_End(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_blankNodePropertyList.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	} else if r0.Rune == ']' {
		r.commit(r0.AsDecodedRunes())

		return readerStack{}, nil
	}

	return readerStack{}, grammar.R_blankNodePropertyList.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
}
