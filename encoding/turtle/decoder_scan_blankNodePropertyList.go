package turtle

import (
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

func reader_scan_blankNodePropertyList_End(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_blankNodePropertyList.Err(r.newOffsetError(err, nil, nil))
	} else if r0 == ']' {
		r.commit([]rune{r0})

		return readerStack{}, nil
	}

	return readerStack{}, grammar.R_blankNodePropertyList.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
}
