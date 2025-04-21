package turtle

import "github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"

func reader_scan_ObjectList_Continue(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_objectList.Err(r.newOffsetError(err, nil, nil))
	}

	if r0 == ',' {
		r.commit([]rune{r0})

		r.pushState(ectx, reader_scan_ObjectList_Continue)

		return readerStack{ectx, reader_scan_Object}, nil
	}

	r.buf.BacktrackRunes(r0)

	return readerStack{}, nil
}
