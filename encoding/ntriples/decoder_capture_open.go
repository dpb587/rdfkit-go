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

func (r *Decoder) captureOpenBlankNode(uncommitted []rune) (rdf.BlankNode, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 2 && uncommitted[0:2] == "_:")

	r0, err := r.buf.NextRune()
	if err != nil {
		return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, uncommitted, nil))
	}

	switch {
	case internal.IsRune_PN_CHARS_U(r0), '0' <= r0 && r0 <= '9':
		// valid
	default:
		return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
	}

	uncommitted = append(uncommitted, r0)

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case internal.IsRune_PN_CHARS(r0), r0 == '.':
			uncommitted = append(uncommitted, r0)
		case unicode.IsSpace(r0):
			r.buf.BacktrackRunes(r0)

			goto DONE
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	if len(uncommitted) > 3 {
		if uncommitted[len(uncommitted)-1] == '.' {
			r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
			uncommitted = uncommitted[0 : len(uncommitted)-1]
		}

		if !internal.IsRune_PN_CHARS(uncommitted[len(uncommitted)-1]) {
			return nil, nil, grammar.R_BLANK_NODE_LABEL.Err(r.newOffsetError(
				cursorioutil.UnexpectedRuneError{
					Rune: uncommitted[len(uncommitted)-1],
				},
				uncommitted[0:len(uncommitted)-1],
				[]rune{uncommitted[len(uncommitted)-1]},
			))
		}
	}

	return r.blankNodeStringMapper.MapBlankNodeIdentifier(string(uncommitted[2:])), r.commitForTextOffsetRange(uncommitted), nil
}

func (r *Decoder) captureOpenIRI(uncommitted []rune) (rdf.IRI, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 1 && uncommitted[0] == '<')

	var decoded []rune

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case r0 == '>':
			uncommitted = append(uncommitted, r0)

			goto DONE
		case r0 == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(err, append(uncommitted, r0), nil))
			}

			switch r1 {
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
				return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
			}
		case 0x00 <= r0 && r0 <= 0x20,
			r0 == '<',
			r0 == '"',
			r0 == '{',
			r0 == '}',
			r0 == '|',
			r0 == '^',
			r0 == '`':
			return "", nil, grammar.R_IRIREF.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
		default:
			decoded = append(decoded, r0)
			uncommitted = append(uncommitted, r0)
		}
	}

DONE:

	urlString := string(decoded)

	cr := r.commitForTextOffsetRange(uncommitted)

	{
		// apparently we should validate these are absolute according to the w3 test suite
		urlParsed, err := url.Parse(urlString)
		if err != nil {
			return "", nil, grammar.R_IRIREF.ErrCursorRange(fmt.Errorf("parse url: %v", err), cr)
		} else if !urlParsed.IsAbs() {
			return "", nil, grammar.R_IRIREF.ErrCursorRange(errors.New("relative urls are not allowed"), cr)
		}
	}

	return rdf.IRI(decoded), cr, nil
}

func (r *Decoder) captureOpenLiteral(uncommitted []rune) (rdf.Literal, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 1 && uncommitted[0] == '"')

	var decoded []rune

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(err, uncommitted, nil)))
		}

		switch {
		case r0 == '"':
			uncommitted = append(uncommitted, r0)

			goto END_LEXICAL_FORM
		case r0 == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(err, append(uncommitted, r0), nil)))
			}

			switch r1 {
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
				return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1})))
			}
		default:
			decoded = append(decoded, r0)
			uncommitted = append(uncommitted, r0)
		}
	}

END_LEXICAL_FORM:

	stringRange := r.commitForTextOffsetRange(uncommitted)

	r0, err := r.buf.NextRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			goto END_TOKEN
		}

		return rdf.Literal{}, nil, grammar.R_literal.Err(grammar.R_STRING_LITERAL_QUOTE.Err(r.newOffsetError(err, uncommitted, nil)))
	}

	switch {
	case r0 == '@':
		langtag, langtagRange, err := r.scanOpenLangtag([]rune{r0})
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
			Tags: map[rdf.LiteralTag]string{
				rdf.LanguageLiteralTag: langtag,
			},
		}, fullRange, nil
	case r0 == '^':
		r1, err := r.buf.NextRune()
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(err, append(uncommitted, r0), nil))
		} else if r1 != '^' {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))
		}

		r.commit([]rune{r0, r1})

		r2, err := r.buf.NextRune()
		if err != nil {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(err, append(uncommitted, r0, r1), nil))
		} else if r2 != '<' {
			return rdf.Literal{}, nil, grammar.R_literal.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, append(uncommitted[:], r0, r1), []rune{r2}))
		}

		iri, iriRange, err := r.captureOpenIRI([]rune{r2})
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

func (r *Decoder) scanOpenLangtag(uncommitted []rune) (string, *cursorio.TextOffsetRange, error) {
	// assert(len(uncommitted) == 1 && uncommitted[0] == '@')

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case 'a' <= r0 && r0 <= 'z', 'A' <= r0 && r0 <= 'Z':
			uncommitted = append(uncommitted, r0)
		case r0 == '-':
			if len(uncommitted) == 1 {
				return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
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
			return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch {
		case 'a' <= r0 && r0 <= 'z', 'A' <= r0 && r0 <= 'Z', '0' <= r0 && r0 <= '9':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto END
		}
	}

END:

	if uncommitted[len(uncommitted)-1] == '-' {
		return "", nil, grammar.R_LANGTAG.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1],
			},
			uncommitted[0:len(uncommitted)-1],
			[]rune{uncommitted[len(uncommitted)-1]},
		))
	}

	r.commit(uncommitted[0:1])

	return string(uncommitted[1:]), r.commitForTextOffsetRange(uncommitted[1:]), nil
}
