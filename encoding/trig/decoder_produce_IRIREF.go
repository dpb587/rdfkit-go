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
func (r *Decoder) produceIRIREF(r0 rune) (*tokenIRIREF, error) {
	if r0 != '<' {
		return nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	var uncommitted = []rune{r0}
	var decoded []rune

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_IRIREF.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case r0 == '>':
			uncommitted = append(uncommitted, r0)

			goto DONE
		case r0 == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_IRIREF.Err(r.newOffsetError(err, append(uncommitted, r0), nil))
			}

			switch r1 {
			case 'u':
				decodedRune, nextUncommitted, err := r.decodeUCHAR4(append(uncommitted, r0, r1))
				if err != nil {
					return nil, grammar.R_IRIREF.Err(err)
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			case 'U':
				decodedRune, nextUncommitted, err := r.decodeUCHAR8(append(uncommitted, r0, r1))
				if err != nil {
					return nil, grammar.R_IRIREF.Err(err)
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			default:
				return nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
			}
		case 0x00 <= r0 && r0 <= 0x20,
			r0 == '<',
			r0 == '"',
			r0 == '{',
			r0 == '}',
			r0 == '|',
			r0 == '^',
			r0 == '`':
			return nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
		default:
			decoded = append(decoded, r0)
			uncommitted = append(uncommitted, r0)
		}
	}

DONE:

	return &tokenIRIREF{
		Offsets: r.commitForTextOffsetRange(uncommitted),
		Decoded: string(decoded),
	}, nil
}
