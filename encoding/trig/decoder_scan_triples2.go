package trig

import (
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func reader_triples2_blankNodePropertyList(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_subject.Err(err)
	} else if r0 == ']' {
		r.commit([]rune{r0})

		r.pushState(ectx, reader_scan_triples_End)
		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

		return readerStack{ectx, reader_scan_PredicateObjectList}, nil
	}

	r.buf.BacktrackRunes(r0)

	r.pushState(ectx, reader_scan_triples_End)
	r.pushState(ectx, reader_scan_PredicateObjectList)
	r.pushState(ectx, reader_scan_blankNodePropertyList_End)
	r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

	return readerStack{ectx, reader_scan_PredicateObjectList}, nil
}
