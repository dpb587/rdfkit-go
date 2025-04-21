package trig

import (
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func reader_scan_wrappedGraph(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_wrappedGraph.Err(err)
	} else if r0 != '{' {
		r.buf.BacktrackRunes(r0)

		return readerStack{}, grammar.R_wrappedGraph.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	r.commit([]rune{r0})

	r.pushState(ectx, reader_scan_wrappedGraph_End)

	return readerStack{ectx, reader_scan_triplesBlock}, nil
}

func reader_scan_wrappedGraph_End(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_wrappedGraph.Err(grammar.R_triplesBlock.Err(err))
	} else if r0 != '}' {
		r.buf.BacktrackRunes(r0)

		return readerStack{}, grammar.R_wrappedGraph.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	r.commit([]rune{r0})

	return readerStack{}, nil
}
