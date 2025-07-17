package trig

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
)

type tokenPNAME_NS struct {
	Offsets       *cursorio.TextOffsetRange
	DecodedString string
}

func (r *Decoder) producePNAME_NS(r0 cursorio.DecodedRune) (*tokenPNAME_NS, error) {
	var uncommitted cursorio.DecodedRuneList
	var namespace []rune

	switch {
	case r0.Rune == ':':
		uncommitted = append(uncommitted, r0)

		goto DONE
	case internal.IsRune_PN_CHARS_BASE(r0.Rune):
		namespace = append(namespace, r0.Rune)
		uncommitted = append(uncommitted, r0)
	default:
		return nil, grammar.R_PNAME_NS.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
	}

	for {
		r1, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_PNAME_NS.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case r1.Rune == ':':
			uncommitted = append(uncommitted, r1)

			goto PN_PREFIX_DONE
		case internal.IsRune_PN_CHARS(r1.Rune), r1.Rune == '.':
			namespace = append(namespace, r1.Rune)
			uncommitted = append(uncommitted, r1)
		default:
			return nil, r.newOffsetError(
				cursorioutil.UnexpectedRuneError{Rune: r1.Rune},
				uncommitted.AsDecodedRunes(),
				r1.AsDecodedRunes(),
			)
		}
	}

PN_PREFIX_DONE:

	if len(uncommitted) > 1 && uncommitted[len(uncommitted)-1].Rune == '.' {
		return nil, grammar.R_PNAME_NS.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1].Rune,
			},
			uncommitted[:len(uncommitted)-1].AsDecodedRunes(),
			uncommitted[len(uncommitted)-1].AsDecodedRunes(),
		))
	}

DONE:

	return &tokenPNAME_NS{
		Offsets:       r.commitForTextOffsetRange(uncommitted.AsDecodedRunes()),
		DecodedString: string(namespace),
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
func (r *Decoder) producePrefixedName(r0 cursorio.DecodedRune) (*tokenPrefixedName, error) {
	var uncommitted cursorio.DecodedRuneList

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

			return nil, grammar.R_PrefixedName.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
		}

		switch {
		case internal.IsRune_PN_CHARS_U(r0.Rune),
			r0.Rune == ':',
			'0' <= r0.Rune && r0.Rune <= '9':
			decodedLocal = append(decodedLocal, r0.Rune)
			uncommitted = append(uncommitted, r0)
		case r0.Rune == '%':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{}))))
			} else if _, ok := internal.HexDecode(r1.Rune); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes())))))
			}

			r2, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0, r1).AsDecodedRunes(), cursorio.DecodedRunes{}))))
			} else if _, ok := internal.HexDecode(r2.Rune); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted[:], r0, r1).AsDecodedRunes(), r2.AsDecodedRunes())))))
			}

			decodedLocal = append(decodedLocal, r0.Rune, r1.Rune, r2.Rune)
			uncommitted = append(uncommitted, r0, r1, r2)
		case r0.Rune == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{})))
			}

			switch r1.Rune {
			case '_', '~', '.', '-', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '/', '?', '#', '@', '%':
				decodedLocal = append(decodedLocal, r1.Rune)
				uncommitted = append(uncommitted, r0, r1)
			default:
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PN_LOCAL_ESC.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes()))))
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

			return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{})))
		}

		switch {
		case internal.IsRune_PN_CHARS(r0.Rune),
			r0.Rune == '.',
			r0.Rune == ':':
			decodedLocal = append(decodedLocal, r0.Rune)
			uncommitted = append(uncommitted, r0)
		case r0.Rune == '%':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{}))))
			} else if _, ok := internal.HexDecode(r1.Rune); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes())))))
			}

			r2, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(r.newOffsetError(err, append(uncommitted, r0, r1).AsDecodedRunes(), cursorio.DecodedRunes{}))))
			} else if _, ok := internal.HexDecode(r2.Rune); !ok {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PERCENT.Err(grammar.R_HEX.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted[:], r0, r1).AsDecodedRunes(), r2.AsDecodedRunes())))))
			}

			decodedLocal = append(decodedLocal, r0.Rune, r1.Rune, r2.Rune)
			uncommitted = append(uncommitted, r0, r1, r2)
		case r0.Rune == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{})))
			}

			switch r1.Rune {
			case '_', '~', '.', '-', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '/', '?', '#', '@', '%':
				decodedLocal = append(decodedLocal, r1.Rune)
				uncommitted = append(uncommitted, r0, r1)
			default:
				return nil, grammar.R_PrefixedName.Err(grammar.R_PN_LOCAL.Err(grammar.R_PN_LOCAL_ESC.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes()))))
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

	cr := r.commitForTextOffsetRange(uncommitted.AsDecodedRunes())

	if cr != nil {
		cr = &cursorio.TextOffsetRange{
			From:  namespaceToken.Offsets.From,
			Until: cr.Until,
		}
	}

	return &tokenPrefixedName{
		Offsets:          cr,
		NamespaceDecoded: namespaceToken.DecodedString,
		LocalDecoded:     string(decodedLocal),
	}, nil
}
