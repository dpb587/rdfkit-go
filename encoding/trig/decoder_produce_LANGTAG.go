package trig

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

type tokenLANGTAG struct {
	Offsets *cursorio.TextOffsetRange
	Decoded string
}

// LANGTAG ::= '@' [a-zA-Z]+ ('-' [a-zA-Z0-9]+)*
func (r *Decoder) produceLANGTAG(r0 rune) (*tokenLANGTAG, error) {
	if r0 != '@' {
		return nil, grammar.R_LANGTAG.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	var uncommitted = []rune{r0}

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(uncommitted) > 1 {
					goto DONE
				}
			}

			return nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case 'a' <= r0 && r0 <= 'z', 'A' <= r0 && r0 <= 'Z':
			uncommitted = append(uncommitted, r0)
		case r0 == '-':
			if len(uncommitted) == 1 {
				return nil, grammar.R_LANGTAG.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
			}

			uncommitted = append(uncommitted, r0)

			goto PRIMARY_DELIMITER_DONE
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

PRIMARY_DELIMITER_DONE:

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case 'a' <= r0 && r0 <= 'z', 'A' <= r0 && r0 <= 'Z', '0' <= r0 && r0 <= '9':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	if uncommitted[len(uncommitted)-1] == '-' {
		return nil, grammar.R_LANGTAG.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1],
			},
			uncommitted[:len(uncommitted)-1],
			[]rune{uncommitted[len(uncommitted)-1]},
		))
	}

	r.commit(uncommitted[0:1])

	valueUncommitted := uncommitted[1:]

	return &tokenLANGTAG{
		Offsets: r.commitForTextOffsetRange(valueUncommitted),
		Decoded: string(valueUncommitted),
	}, nil
}
