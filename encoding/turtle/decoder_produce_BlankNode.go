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

func (r *Decoder) produceBlankNode(r0 rune) (*tokenBlankNode, error) {
	if r0 != '_' {
		return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	var uncommitted = []rune{}

	{
		r1, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
		} else if r1 != ':' {
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, append(uncommitted, r0, r1), nil))
		}

		switch {
		case internal.IsRune_PN_CHARS_U(r2), '0' <= r2 && r2 <= '9':
			// valid
		default:
			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, append(uncommitted[:], r0, r1), []rune{r2}))
		}

		r.commit([]rune{r0, r1})

		uncommitted = append(uncommitted, r2)
	}

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case internal.IsRune_PN_CHARS(r0), r0 == '.':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	if uncommitted[len(uncommitted)-1] == '.' {
		r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
		uncommitted = uncommitted[0 : len(uncommitted)-1]
	}

	if len(uncommitted) > 1 && !internal.IsRune_PN_CHARS(uncommitted[len(uncommitted)-1]) {
		return nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1],
			},
			uncommitted[0:len(uncommitted)-1],
			[]rune{uncommitted[len(uncommitted)-1]},
		))
	}

	token := &tokenBlankNode{
		Offsets: r.commitForTextOffsetRange(uncommitted),
		Decoded: string(uncommitted),
	}

	return token, nil
}
