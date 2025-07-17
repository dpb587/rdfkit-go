package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func reader_scan_wrappedGraph(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_wrappedGraph.Err(err)
	} else if r0.Rune != '{' {
		r.buf.BacktrackRunes(r0)

		return readerStack{}, grammar.R_wrappedGraph.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
	}

	r.commit(r0.AsDecodedRunes())

	r.pushState(ectx, reader_scan_wrappedGraph_End)

	return readerStack{ectx, reader_scan_triplesBlock}, nil
}

func reader_scan_wrappedGraph_End(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_wrappedGraph.Err(grammar.R_triplesBlock.Err(err))
	} else if r0.Rune != '}' {
		r.buf.BacktrackRunes(r0)

		return readerStack{}, grammar.R_wrappedGraph.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
	}

	r.commit(r0.AsDecodedRunes())

	return readerStack{}, nil
}
