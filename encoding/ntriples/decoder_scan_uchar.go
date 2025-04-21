package ntriples

import (
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal/grammar"
)

func (r *Decoder) decodeUCHAR4(uncommitted []rune) (rune, []rune, error) {
	r0, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r0x, ok := internal.HexDecode(r0)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
	}

	r1, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r1x, ok := internal.HexDecode(r1)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
	}

	r2, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r2x, ok := internal.HexDecode(r2)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, append(uncommitted[:], r0, r1), []rune{r2}))
	}

	r3, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r3x, ok := internal.HexDecode(r3)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3}, append(uncommitted[:], r0, r1, r2), []rune{r3}))
	}

	return rune(r0x<<12 | r1x<<8 | r2x<<4 | r3x),
		append(uncommitted, r0, r1, r2, r3),
		nil
}

func (r *Decoder) decodeUCHAR8(uncommitted []rune) (rune, []rune, error) {
	r0, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r0x, ok := internal.HexDecode(r0)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
	} else if r0x > 0 {
		return 0, nil, grammar.R_UCHAR.Err(encoding.ExceedsMaxUnicodePointErr)
	}

	r1, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r1x, ok := internal.HexDecode(r1)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
	} else if r1x > 0 {
		return 0, nil, grammar.R_UCHAR.Err(encoding.ExceedsMaxUnicodePointErr)
	}

	r2, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r2x, ok := internal.HexDecode(r2)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, append(uncommitted[:], r0, r1), []rune{r2}))
	} else if r2x > 1 {
		return 0, nil, grammar.R_UCHAR.Err(encoding.ExceedsMaxUnicodePointErr)
	}

	r3, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r3x, ok := internal.HexDecode(r3)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3}, append(uncommitted[:], r0, r1, r2), []rune{r3}))
	}

	r4, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r4x, ok := internal.HexDecode(r4)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r4}, append(uncommitted[:], r0, r1, r2, r3), []rune{r4}))
	}

	r5, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r5x, ok := internal.HexDecode(r5)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r5}, append(uncommitted[:], r0, r1, r2, r3, r4), []rune{r5}))
	}

	r6, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r6x, ok := internal.HexDecode(r6)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r6}, append(uncommitted[:], r0, r1, r2, r3, r4, r5), []rune{r6}))
	}

	r7, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r7x, ok := internal.HexDecode(r7)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r7}, append(uncommitted[:], r0, r1, r2, r3, r4, r5, r6), []rune{r7}))
	}

	return rune(r0x<<28 | r1x<<24 | r2x<<20 | r3x<<16 | r4x<<12 | r5x<<8 | r6x<<4 | r7x),
		[]rune{r0, r1, r2, r3, r4, r5, r6, r7},
		nil
}
