package turtle

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

type tokenLANGTAG struct {
	Offsets *cursorio.TextOffsetRange
	Decoded string
}

// LANGTAG ::= '@' [a-zA-Z]+ ('-' [a-zA-Z0-9]+)*
func (r *Decoder) produceLANGTAG(r0 cursorio.DecodedRune) (*tokenLANGTAG, error) {
	if r0.Rune != '@' {
		return nil, grammar.R_LANGTAG.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
	}

	var uncommitted cursorio.DecodedRuneList = cursorio.DecodedRuneList{r0}

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(uncommitted) > 1 {
					goto DONE
				}
			}

			return nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case 'a' <= r0.Rune && r0.Rune <= 'z', 'A' <= r0.Rune && r0.Rune <= 'Z':
			uncommitted = append(uncommitted, r0)
		case r0.Rune == '-':
			if len(uncommitted) == 1 {
				return nil, grammar.R_LANGTAG.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, uncommitted.AsDecodedRunes(), r0.AsDecodedRunes()))
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

			return nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case 'a' <= r0.Rune && r0.Rune <= 'z', 'A' <= r0.Rune && r0.Rune <= 'Z', '0' <= r0.Rune && r0.Rune <= '9':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	if uncommitted[len(uncommitted)-1].Rune == '-' {
		return nil, grammar.R_LANGTAG.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1].Rune,
			},
			uncommitted[:len(uncommitted)-1].AsDecodedRunes(),
			uncommitted[len(uncommitted)-1].AsDecodedRunes(),
		))
	}

	r.commit(uncommitted[0:1].AsDecodedRunes())

	valueUncommitted := uncommitted[1:]

	return &tokenLANGTAG{
		Offsets: r.commitForTextOffsetRange(valueUncommitted.AsDecodedRunes()),
		Decoded: valueUncommitted.AsDecodedRunes().String(),
	}, nil
}
