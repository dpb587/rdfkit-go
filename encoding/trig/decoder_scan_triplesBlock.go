package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func reader_scan_triplesBlock(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triplesBlock.Err(err)
	}

	switch r0.Rune {
	case '}':
		r.buf.BacktrackRunes(r0)

		return readerStack{}, nil
	}

	r.pushState(ectx, reader_scan_triplesBlock_QUEST)

	r.buf.BacktrackRunes(r0)

	return readerStack{ectx, reader_scan_triples}, nil
}

func reader_scan_triplesBlock_QUEST(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_wrappedGraph.Err(grammar.R_triplesBlock.Err(err))
	} else if r0.Rune == '.' {
		r.commit(r0.AsDecodedRunes())

		return readerStack{ectx, reader_scan_triplesBlock}, nil
	} else if r0.Rune == '}' {
		r.buf.BacktrackRunes(r0)

		return readerStack{}, nil
	}

	r.buf.BacktrackRunes(r0)

	return readerStack{ectx, reader_scan_triplesBlock}, nil
}
