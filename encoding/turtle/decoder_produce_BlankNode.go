package turtle

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

type tokenBlankNode struct {
	Offsets *cursorio.TextOffsetRange
	Decoded string
}

func (r *Decoder) produceBlankNode(r0 cursorio.DecodedRune) (*tokenBlankNode, error) {
	if r0.Rune != '_' {
		return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
	}

	var uncommitted cursorio.DecodedRuneList

	{
		r1, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted, r0).AsDecodedRunes(), r1.AsDecodedRunes()))
		} else if r1.Rune != ':' {
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted, r0).AsDecodedRunes(), r1.AsDecodedRunes()))
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, append(uncommitted, r0, r1).AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case internal.IsRune_PN_CHARS_U(r2.Rune), '0' <= r2.Rune && r2.Rune <= '9':
			// valid
		default:
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted, r0, r1).AsDecodedRunes(), r2.AsDecodedRunes()))
		}

		r.commit(cursorio.NewDecodedRunes(r0, r1))

		uncommitted = append(uncommitted, r2)
	}

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case internal.IsRune_PN_CHARS(r0.Rune), r0.Rune == '.':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	if uncommitted[len(uncommitted)-1].Rune == '.' {
		r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
		uncommitted = uncommitted[0 : len(uncommitted)-1]
	}

	if len(uncommitted) > 1 && !internal.IsRune_PN_CHARS(uncommitted[len(uncommitted)-1].Rune) {
		return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1].Rune,
			},
			uncommitted[0:len(uncommitted)-1].AsDecodedRunes(),
			uncommitted[len(uncommitted)-1].AsDecodedRunes(),
		))
	}

	token := &tokenBlankNode{
		Offsets: r.commitForTextOffsetRange(uncommitted.AsDecodedRunes()),
		Decoded: uncommitted.AsDecodedRunes().String(),
	}

	return token, nil
}
