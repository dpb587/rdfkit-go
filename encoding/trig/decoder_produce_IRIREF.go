package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

type tokenIRIREF struct {
	Offsets *cursorio.TextOffsetRange
	Decoded string
}

// IRIREF ::= '<' ([^#x00-#x20<>"{}|^`\] | UCHAR)* '>'
func (r *Decoder) produceIRIREF(r0 cursorio.DecodedRune) (*tokenIRIREF, error) {
	if r0.Rune != '<' {
		return nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
	}

	var uncommitted = cursorio.DecodedRuneList{r0}
	var decoded []rune

	for {
		r1, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_IRIREF.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case r1.Rune == '>':
			uncommitted = append(uncommitted, r1)

			goto DONE
		case r1.Rune == '\\':
			r2, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_IRIREF.Err(r.newOffsetError(err, append(uncommitted, r1).AsDecodedRunes(), cursorio.DecodedRunes{}))
			}

			switch r2.Rune {
			case 'u':
				decodedRune, nextUncommitted, err := r.decodeUCHAR4(append(uncommitted, r1, r2))
				if err != nil {
					return nil, grammar.R_IRIREF.Err(err)
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			case 'U':
				decodedRune, nextUncommitted, err := r.decodeUCHAR8(append(uncommitted, r1, r2))
				if err != nil {
					return nil, grammar.R_IRIREF.Err(err)
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			default:
				return nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted[:], r1).AsDecodedRunes(), r2.AsDecodedRunes()))
			}
		case 0x00 <= r1.Rune && r1.Rune <= 0x20,
			r1.Rune == '<',
			r1.Rune == '"',
			r1.Rune == '{',
			r1.Rune == '}',
			r1.Rune == '|',
			r1.Rune == '^',
			r1.Rune == '`':
			return nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, uncommitted.AsDecodedRunes(), r1.AsDecodedRunes()))
		default:
			decoded = append(decoded, r1.Rune)
			uncommitted = append(uncommitted, r1)
		}
	}

DONE:

	return &tokenIRIREF{
		Offsets: r.commitForTextOffsetRange(uncommitted.AsDecodedRunes()),
		Decoded: string(decoded),
	}, nil
}
