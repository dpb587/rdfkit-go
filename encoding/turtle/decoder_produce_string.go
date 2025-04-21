package turtle

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

type tokenString struct {
	Offsets *cursorio.TextOffsetRange
	Decoded string
}

// String                           ::= STRING_LITERAL_QUOTE | STRING_LITERAL_SINGLE_QUOTE | STRING_LITERAL_LONG_SINGLE_QUOTE | STRING_LITERAL_LONG_QUOTE
// STRING_LITERAL_QUOTE             ::= '"' ([^#x22#x5C#xA#xD] | ECHAR | UCHAR)* '"'
// STRING_LITERAL_SINGLE_QUOTE      ::= "'" ([^#x27#x5C#xA#xD] | ECHAR | UCHAR)* "'"
// STRING_LITERAL_LONG_SINGLE_QUOTE ::= "”'" (("'" | "”")? ([^'\] | ECHAR | UCHAR))* "”'"
// STRING_LITERAL_LONG_QUOTE        ::= '"""' (('"' | '""')? ([^"\] | ECHAR | UCHAR))* '"""'
func (r *Decoder) produceString(r0 rune) (*tokenString, error) {
	var grammarRule grammar.R

	switch r0 {
	case '"':
		grammarRule = grammar.R_STRING_LITERAL_QUOTE
	case '\'':
		grammarRule = grammar.R_STRING_LITERAL_SINGLE_QUOTE
	default:
		return nil, grammar.R_String.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
	}

	var uncommitted = []rune{r0}
	var delimiterRune = r0
	var delimiterTriple = false
	var decoded []rune

	r0, err := r.buf.NextRune()
	if err != nil {
		return nil, grammar.R_String.Err(grammarRule.Err(r.newOffsetError(err, uncommitted, nil)))
	}

	if r0 == delimiterRune {
		r1, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammar.R_String.Err(grammarRule.Err(r.newOffsetError(err, append(uncommitted, r0), nil)))
		} else if r1 == delimiterRune {
			delimiterTriple = true

			switch grammarRule {
			case grammar.R_STRING_LITERAL_QUOTE:
				grammarRule = grammar.R_STRING_LITERAL_LONG_QUOTE
			case grammar.R_STRING_LITERAL_SINGLE_QUOTE:
				grammarRule = grammar.R_STRING_LITERAL_LONG_SINGLE_QUOTE
			}

			uncommitted = append(uncommitted, r0, r1)

			goto START_DELIMITER_DONE
		}

		r.buf.BacktrackRunes(r1)

		return &tokenString{
			Offsets: r.commitForTextOffsetRange(append(uncommitted, r0, r1)),
			Decoded: "",
		}, nil
	} else {
		r.buf.BacktrackRunes(r0)
	}

START_DELIMITER_DONE:

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, grammar.R_String.Err(grammarRule.Err(r.newOffsetError(err, uncommitted, nil)))
		}

		switch {
		case r0 == '"' || r0 == '\'':
			if r0 == delimiterRune {
				if !delimiterTriple {
					uncommitted = append(uncommitted, r0)

					goto DONE
				}

				r1, err := r.buf.NextRune()
				if err != nil {
					return nil, grammar.R_String.Err(grammarRule.Err(r.newOffsetError(err, append(uncommitted, r0), nil)))
				} else if r1 == delimiterRune {
					r2, err := r.buf.NextRune()
					if err != nil {
						return nil, grammar.R_String.Err(grammarRule.Err(r.newOffsetError(err, append(uncommitted, r0, r1), nil)))
					} else if r2 == delimiterRune {
						uncommitted = append(uncommitted, r0, r1, r2)

						goto DONE
					}

					r.buf.BacktrackRunes(r1, r2)
				} else {
					r.buf.BacktrackRunes(r1)
				}
			}

			decoded = append(decoded, r0)
			uncommitted = append(uncommitted, r0)
		case r0 == '\\':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, grammar.R_String.Err(grammarRule.Err(r.newOffsetError(err, uncommitted, nil)))
			}

			switch r1 {
			case 'u':
				decodedRune, nextUncommitted, err := r.decodeUCHAR4(append(uncommitted, r0, r1))
				if err != nil {
					return nil, grammar.R_String.Err(grammarRule.Err(err))
				}

				decoded = append(decoded, decodedRune)
				uncommitted = nextUncommitted
			case 'U':
				decodedRune, nextUncommitted, err := r.decodeUCHAR8(append(uncommitted, r0, r1))
				if err != nil {
					return nil, grammar.R_String.Err(grammarRule.Err(err))
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
				return nil, grammar.R_String.Err(grammarRule.Err(grammar.R_ECHAR.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, append(uncommitted[:], r0), []rune{r1}))))
			}
		default:
			decoded = append(decoded, r0)
			uncommitted = append(uncommitted, r0)
		}
	}

DONE:

	return &tokenString{
		Offsets: r.commitForTextOffsetRange(uncommitted),
		Decoded: string(decoded),
	}, nil
}
