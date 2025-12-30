package ntriples

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal/grammar"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func (r *Decoder) captureOpenBlankNode(uncommitted cursorio.DecodedRuneList) (rdf.BlankNode, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 2 && uncommitted[0:2] == "_:")

	r0, err := r.buf.NextRune()
	if err != nil {
		return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
	}

	switch {
	case internal.IsRune_PN_CHARS_U(r0.Rune), '0' <= r0.Rune && r0.Rune <= '9':
		// valid
	default:
		return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, uncommitted.AsDecodedRunes(), r0.AsDecodedRunes()))
	}

	uncommitted = append(uncommitted, r0)

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case internal.IsRune_PN_CHARS(r0.Rune), r0.Rune == '.':
			uncommitted = append(uncommitted, r0)
		case unicode.IsSpace(r0.Rune):
			r.buf.BacktrackRunes(r0)

			goto DONE
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	if len(uncommitted) > 3 {
		if uncommitted[len(uncommitted)-1].Rune == '.' {
			r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
			uncommitted = uncommitted[0 : len(uncommitted)-1]
		}

		if !internal.IsRune_PN_CHARS(uncommitted[len(uncommitted)-1].Rune) {
			return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(
				cursorioutil.UnexpectedRuneError{
					Rune: uncommitted[len(uncommitted)-1].Rune,
				},
				uncommitted[0:len(uncommitted)-1].AsDecodedRunes(),
				uncommitted[len(uncommitted)-1].AsDecodedRunes(),
			))
		}
	}

	return r.blankNodeStringMapper.MapBlankNodeIdentifier(string(uncommitted[2:].AsDecodedRunes().Runes)), r.commitForTextOffsetRange(uncommitted.AsDecodedRunes()), nil
}

func (r *Decoder) captureOpenIRI(uncommitted cursorio.DecodedRuneList) (rdf.IRI, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 1 && uncommitted[0] == '<')

	var decoded []rune

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case r0.Rune == '>':
			uncommitted = append(uncommitted, r0)

			goto DONE
		case r0.Rune == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{}))
			}

			switch r1.Rune {
			case 'u':
				decodedRune, nextUncommitted, err := r.decodeUCHAR4(append(uncommitted, r0, r1))
				if err != nil {
					return "", nil, grammar.R_IRIREF.Err(err)
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			case 'U':
				decodedRune, nextUncommitted, err := r.decodeUCHAR8(append(uncommitted, r0, r1))
				if err != nil {
					return "", nil, grammar.R_IRIREF.Err(err)
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			default:
				return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes()))
			}
		case 0x00 <= r0.Rune && r0.Rune <= 0x20,
			r0.Rune == '<',
			r0.Rune == '"',
			r0.Rune == '{',
			r0.Rune == '}',
			r0.Rune == '|',
			r0.Rune == '^',
			r0.Rune == '`':
			return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, uncommitted.AsDecodedRunes(), r0.AsDecodedRunes()))
		default:
			decoded = append(decoded, r0.Rune)
			uncommitted = append(uncommitted, r0)
		}
	}

DONE:

	urlString := string(decoded)

	cr := r.commitForTextOffsetRange(uncommitted.AsDecodedRunes())

	{
		// apparently we should validate these are absolute according to the w3 test suite
		urlParsed, err := url.Parse(urlString)
		if err != nil {
			return "", nil, grammar.R_IRIREF.ErrWithTextOffsetRange(fmt.Errorf("parse url: %v", err), cr)
		} else if !urlParsed.IsAbs() {
			return "", nil, grammar.R_IRIREF.ErrWithTextOffsetRange(errors.New("relative urls are not allowed"), cr)
		}
	}

	return rdf.IRI(decoded), cr, nil
}

func (r *Decoder) captureOpenLiteral(uncommitted cursorio.DecodedRuneList) (rdf.Literal, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 1 && uncommitted[0] == '"')

	var decoded []rune

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{})))
		}

		switch {
		case r0.Rune == '"':
			uncommitted = append(uncommitted, r0)

			goto END_LEXICAL_FORM
		case r0.Rune == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{})))
			}

			switch r1.Rune {
			case 'u':
				decodedRune, nextUncommitted, err := r.decodeUCHAR4(append(uncommitted, r0, r1))
				if err != nil {
					return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(err))
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			case 'U':
				decodedRune, nextUncommitted, err := r.decodeUCHAR8(append(uncommitted, r0, r1))
				if err != nil {
					return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(err))
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			case 't':
				decoded = append(decoded, '\t')
				uncommitted = append(uncommitted, r0, r1)
			case 'b':
				decoded = append(decoded, '\b')
				uncommitted = append(uncommitted, r0, r1)
			case 'n':
				decoded = append(decoded, '\n')
				uncommitted = append(uncommitted, r0, r1)
			case 'r':
				decoded = append(decoded, '\r')
				uncommitted = append(uncommitted, r0, r1)
			case 'f':
				decoded = append(decoded, '\f')
				uncommitted = append(uncommitted, r0, r1)
			case '"':
				decoded = append(decoded, '"')
				uncommitted = append(uncommitted, r0, r1)
			case '\'':
				decoded = append(decoded, '\'')
				uncommitted = append(uncommitted, r0, r1)
			case '\\':
				decoded = append(decoded, '\\')
				uncommitted = append(uncommitted, r0, r1)
			default:
				return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes())))
			}
		default:
			decoded = append(decoded, r0.Rune)
			uncommitted = append(uncommitted, r0)
		}
	}

END_LEXICAL_FORM:

	stringRange := r.commitForTextOffsetRange(uncommitted.AsDecodedRunes())

	r0, err := r.buf.NextRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			goto END_TOKEN
		}

		return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{})))
	}

	switch {
	case r0.Rune == '@':
		langtag, langtagRange, err := r.scanOpenLangtag(cursorio.DecodedRuneList{r0})
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(err)
		}

		var fullRange *cursorio.TextOffsetRange

		if stringRange != nil && langtagRange != nil {
			fullRange = &cursorio.TextOffsetRange{
				From:  stringRange.From,
				Until: langtagRange.Until,
			}
		}

		return rdf.Literal{
			Datatype:    rdfiri.LangString_Datatype,
			LexicalForm: string(decoded),
			Tag: rdf.LanguageLiteralTag{
				Language: langtag,
			},
		}, fullRange, nil
	case r0.Rune == '^':
		r1, err := r.buf.NextRune()
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(err, append(uncommitted, r0).AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r1.Rune != '^' {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, append(uncommitted[:], r0).AsDecodedRunes(), r1.AsDecodedRunes()))
		}

		r.commit(cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes())

		r2, err := r.buf.NextRune()
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(err, append(uncommitted, r0, r1).AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r2.Rune != '<' {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, append(uncommitted[:], r0, r1).AsDecodedRunes(), r2.AsDecodedRunes()))
		}

		iri, iriRange, err := r.captureOpenIRI(cursorio.DecodedRuneList{r2})
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(err)
		}

		var fullRange *cursorio.TextOffsetRange

		if stringRange != nil && iriRange != nil {
			fullRange = &cursorio.TextOffsetRange{
				From:  stringRange.From,
				Until: iriRange.Until,
			}
		}

		return rdf.Literal{
			Datatype:    iri,
			LexicalForm: string(decoded),
		}, fullRange, nil
	default:
		r.buf.BacktrackRunes(r0)
	}

END_TOKEN:

	return rdf.Literal{
		Datatype:    xsdiri.String_Datatype,
		LexicalForm: string(decoded),
	}, stringRange, nil
}

func (r *Decoder) scanOpenLangtag(uncommitted cursorio.DecodedRuneList) (string, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 1 && uncommitted[0] == '@')

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case 'a' <= r0.Rune && r0.Rune <= 'z', 'A' <= r0.Rune && r0.Rune <= 'Z':
			uncommitted = append(uncommitted, r0)
		case r0.Rune == '-':
			if len(uncommitted) == 1 {
				return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, uncommitted.AsDecodedRunes(), r0.AsDecodedRunes()))
			}

			uncommitted = append(uncommitted, r0)

			goto SECONDARY
		default:
			r.buf.BacktrackRunes(r0)

			goto END
		}
	}

SECONDARY:

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		switch {
		case 'a' <= r0.Rune && r0.Rune <= 'z', 'A' <= r0.Rune && r0.Rune <= 'Z', '0' <= r0.Rune && r0.Rune <= '9':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto END
		}
	}

END:

	if uncommitted[len(uncommitted)-1].Rune == '-' {
		return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1].Rune,
			},
			uncommitted[0:len(uncommitted)-1].AsDecodedRunes(),
			uncommitted[len(uncommitted)-1].AsDecodedRunes(),
		))
	}

	r.commit(uncommitted[0:1].AsDecodedRunes())

	return string(uncommitted[1:].AsDecodedRunes().Runes), r.commitForTextOffsetRange(uncommitted[1:].AsDecodedRunes()), nil
}
