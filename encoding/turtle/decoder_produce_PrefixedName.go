package turtle

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

type tokenPNAME_NS struct {
	Offsets *cursorio.TextOffsetRange
	Decoded string
}

func (r *Decoder) producePNAME_NS(r0 rune) (*tokenPNAME_NS, error) {
	var uncommitted []rune
	var namespace []rune

	switch {
	case r0 == ':':
		uncommitted = append(uncommitted, r0)

		goto DONE
	case internal.IsRune_PN_CHARS_BASE(r0):
		namespace = append(namespace, r0)
		uncommitted = append(uncommitted, r0)
	default:
		return nil, grammar.R_PNAME_NS.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_PNAME_NS.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case r0 == ':':
			uncommitted = append(uncommitted, r0)

			goto PN_PREFIX_DONE
		case internal.IsRune_PN_CHARS(r0), r0 == '.':
			namespace = append(namespace, r0)
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto PN_PREFIX_DONE
		}
	}

PN_PREFIX_DONE:

	if len(uncommitted) > 1 && uncommitted[len(uncommitted)-1] == '.' {
		return nil, grammar.R_PNAME_NS.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1],
			},
			uncommitted[:len(uncommitted)-1],
			[]rune{uncommitted[len(uncommitted)-1]},
		))
	}

DONE:

	return &tokenPNAME_NS{
		Offsets: r.commitForTextOffsetRange(uncommitted),
		Decoded: string(namespace),
	}, nil
}

type tokenPrefixedName struct {
	Offsets          *cursorio.TextOffsetRange
	NamespaceDecoded string
	LocalDecoded     string
}

// PrefixedName ::= PNAME_LN | PNAME_NS
// PNAME_NS     ::= PN_PREFIX? ':'
// PNAME_LN     ::= PNAME_NS PN_LOCAL
func (r *Decoder) producePrefixedName(r0 rune) (*tokenPrefixedName, error) {
	var uncommitted []rune

	var decodedLocal []rune

	namespaceToken, err := r.producePNAME_NS(r0)
	if err != nil {
		return nil, grammar.R_PrefixedName.Err(err)
	}

	{
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammar.R_PrefixedName.Err(r.newOffsetError(err, nil, nil))
		}

		switch {
		case internal.IsRune_PN_CHARS_U(r0),
			r0 == ':',
			'0' <= r0 && r0 <= '9':
			decodedLocal = append(decodedLocal, r0)
			uncommitted = append(uncommitted, r0)
		case r0 == '%':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0), nil))))
			} else if _, ok := internal.HexDecode(r1); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1})))))
			}

			r2, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0, r1), nil))))
			} else if _, ok := internal.HexDecode(r2); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, append(uncommitted[:], r0, r1), []rune{r2})))))
			}

			decodedLocal = append(decodedLocal, r0, r1, r2)
			uncommitted = append(uncommitted, r0, r1, r2)
		case r0 == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(r.newOffsetError(err, append(uncommitted, r0), nil)))
			}

			switch r1 {
			case '_', '~', '.', '-', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '/', '?', '#', '@', '%':
				decodedLocal = append(decodedLocal, r1)
				uncommitted = append(uncommitted, r0, r1)
			default:
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PN_LOCAL_ESC.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))))
			}
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto PN_LOCAL_DONE
			}

			return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(r.newOffsetError(err, uncommitted, nil)))
		}

		switch {
		case internal.IsRune_PN_CHARS(r0),
			r0 == '.',
			r0 == ':':
			decodedLocal = append(decodedLocal, r0)
			uncommitted = append(uncommitted, r0)
		case r0 == '%':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0), nil))))
			} else if _, ok := internal.HexDecode(r1); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1})))))
			}

			r2, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0, r1), nil))))
			} else if _, ok := internal.HexDecode(r2); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, append(uncommitted[:], r0, r1), []rune{r2})))))
			}

			decodedLocal = append(decodedLocal, r0, r1, r2)
			uncommitted = append(uncommitted, r0, r1, r2)
		case r0 == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(r.newOffsetError(err, append(uncommitted, r0), nil)))
			}

			switch r1 {
			case '_', '~', '.', '-', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '/', '?', '#', '@', '%':
				decodedLocal = append(decodedLocal, r1)
				uncommitted = append(uncommitted, r0, r1)
			default:
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PN_LOCAL_ESC.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))))
			}
		default:
			r.buf.BacktrackRunes(r0)

			goto PN_LOCAL_DONE
		}
	}

PN_LOCAL_DONE:

	if decodedLocal[len(decodedLocal)-1] == '.' {
		r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
		uncommitted = uncommitted[0 : len(uncommitted)-1]
		decodedLocal = decodedLocal[0 : len(decodedLocal)-1]

		if decodedLocal[len(decodedLocal)-1] == '\\' {
			r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
			uncommitted = uncommitted[0 : len(uncommitted)-1]
			decodedLocal = decodedLocal[0 : len(decodedLocal)-1]
		}
	}

DONE:

	cr := r.commitForTextOffsetRange(uncommitted)

	if cr != nil {
		cr = &cursorio.TextOffsetRange{
			From:  namespaceToken.Offsets.From,
			Until: cr.Until,
		}
	}

	return &tokenPrefixedName{
		Offsets:          cr,
		NamespaceDecoded: namespaceToken.Decoded,
		LocalDecoded:     string(decodedLocal),
	}, nil
}
