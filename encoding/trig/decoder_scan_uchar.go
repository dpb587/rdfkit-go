package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

func (r *Decoder) decodeUCHAR4(uncommitted cursorio.DecodedRuneList) (rune, cursorio.DecodedRuneList, error) {
	r0, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r0x, ok := internal.HexDecode(r0.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, uncommitted.AsDecodedRunes(), r0.AsDecodedRunes()))
	}

	r1, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r1x, ok := internal.HexDecode(r1.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted, r0).AsDecodedRunes(), r1.AsDecodedRunes()))
	}

	r2, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r2x, ok := internal.HexDecode(r2.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted, r0, r1).AsDecodedRunes(), r2.AsDecodedRunes()))
	}

	r3, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r3x, ok := internal.HexDecode(r3.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3.Rune}, append(uncommitted, r0, r1, r2).AsDecodedRunes(), r3.AsDecodedRunes()))
	}

	return rune(r0x<<12 | r1x<<8 | r2x<<4 | r3x),
		append(uncommitted, r0, r1, r2, r3),
		nil
}

func (r *Decoder) decodeUCHAR8(uncommitted cursorio.DecodedRuneList) (rune, cursorio.DecodedRuneList, error) {
	r0, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r0x, ok := internal.HexDecode(r0.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, uncommitted.AsDecodedRunes(), r0.AsDecodedRunes()))
	} else if r0x > 0 {
		return 0, nil, grammar.R_UCHAR.Err(encoding.ExceedsMaxUnicodePointErr)
	}

	r1, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r1x, ok := internal.HexDecode(r1.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted, r0).AsDecodedRunes(), r1.AsDecodedRunes()))
	} else if r1x > 0 {
		return 0, nil, grammar.R_UCHAR.Err(encoding.ExceedsMaxUnicodePointErr)
	}

	r2, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r2x, ok := internal.HexDecode(r2.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted, r0, r1).AsDecodedRunes(), r2.AsDecodedRunes()))
	} else if r2x > 1 {
		return 0, nil, grammar.R_UCHAR.Err(encoding.ExceedsMaxUnicodePointErr)
	}

	r3, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r3x, ok := internal.HexDecode(r3.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3.Rune}, append(uncommitted, r0, r1, r2).AsDecodedRunes(), r3.AsDecodedRunes()))
	}

	r4, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r4x, ok := internal.HexDecode(r4.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r4.Rune}, append(uncommitted, r0, r1, r2, r3).AsDecodedRunes(), r4.AsDecodedRunes()))
	}

	r5, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r5x, ok := internal.HexDecode(r5.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r5.Rune}, append(uncommitted, r0, r1, r2, r3, r4).AsDecodedRunes(), r5.AsDecodedRunes()))
	}

	r6, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r6x, ok := internal.HexDecode(r6.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r6.Rune}, append(uncommitted, r0, r1, r2, r3, r4, r5).AsDecodedRunes(), r6.AsDecodedRunes()))
	}

	r7, err := r.buf.NextRune()
	if err != nil {
		return 0, nil, grammar.R_UCHAR.Err(err)
	}

	r7x, ok := internal.HexDecode(r7.Rune)
	if !ok {
		return 0, nil, grammar.R_UCHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r7.Rune}, append(uncommitted, r0, r1, r2, r3, r4, r5, r6).AsDecodedRunes(), r7.AsDecodedRunes()))
	}

	return rune(r0x<<28 | r1x<<24 | r2x<<20 | r3x<<16 | r4x<<12 | r5x<<8 | r6x<<4 | r7x),
		append(uncommitted, r0, r1, r2, r3, r4, r5, r6, r7),
		nil
}
