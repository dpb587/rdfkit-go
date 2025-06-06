package turtle

import (
	"errors"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

type tokenNumericLiteral struct {
	Offsets     *cursorio.TextOffsetRange
	GrammarRule grammar.R
	Decoded     string
}

// NumericLiteral ::= INTEGER | DECIMAL | DOUBLE
// INTEGER        ::= [+-]? [0-9]+
// DECIMAL        ::= [+-]? [0-9]* '.' [0-9]+
// DOUBLE         ::= [+-]? ([0-9]+ '.' [0-9]* EXPONENT | '.' [0-9]+ EXPONENT | [0-9]+ EXPONENT)
// EXPONENT       ::= [eE] [+-]? [0-9]+
func (r *Decoder) produceNumericLiteral(r0 rune) (*tokenNumericLiteral, error) {
	var uncommitted []rune
	var grammarToken = grammar.R_NumericLiteral

	switch r0 {
	case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		uncommitted = append(uncommitted, r0)

		goto SIGN_DONE
	case '.':
		uncommitted = append(uncommitted, r0)
		grammarToken = grammar.R_DECIMAL

		goto INTEGER_DONE
	default:
		return nil, grammarToken.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0}))
	}

SIGN_DONE:

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammarToken.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch r0 {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			uncommitted = append(uncommitted, r0)
		case '.':
			uncommitted = append(uncommitted, r0)
			grammarToken = grammar.R_DECIMAL

			goto INTEGER_DONE
		case 'e', 'E':
			uncommitted = append(uncommitted, r0)
			grammarToken = grammar.R_DOUBLE

			goto DECIMAL_DONE
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

INTEGER_DONE:

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammarToken.Err(r.newOffsetError(err, uncommitted, nil))
		}

		switch r0 {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			uncommitted = append(uncommitted, r0)
		case 'e', 'E':
			uncommitted = append(uncommitted, r0)
			grammarToken = grammar.R_DOUBLE

			goto DECIMAL_DONE
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DECIMAL_DONE:

	{
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, grammarToken.Err(grammar.R_EXPONENT.Err(r.newOffsetError(err, uncommitted, nil)))
		}

		switch r0 {
		case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			uncommitted = append(uncommitted, r0)

			goto EXPONENT_SIGN_DONE
		default:
			return nil, grammarToken.Err(grammar.R_EXPONENT.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, uncommitted, []rune{r0})))
		}
	}

EXPONENT_SIGN_DONE:

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto DONE
			}

			return nil, grammarToken.Err(grammar.R_EXPONENT.Err(r.newOffsetError(err, uncommitted, nil)))
		}

		switch r0 {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			uncommitted = append(uncommitted, r0)
		default:
			r.buf.BacktrackRunes(r0)

			goto DONE
		}
	}

DONE:

	switch uncommitted[len(uncommitted)-1] {
	case '.':
		r.buf.BacktrackRunes(uncommitted[len(uncommitted)-1])
		uncommitted = uncommitted[:len(uncommitted)-1]

		grammarToken = grammar.R_INTEGER
	case '-', '+', 'e', 'E':
		return nil, grammarToken.Err(r.newOffsetError(
			cursorioutil.UnexpectedRuneError{
				Rune: uncommitted[len(uncommitted)-1],
			},
			uncommitted[:len(uncommitted)-1],
			[]rune{uncommitted[len(uncommitted)-1]},
		))
	}

	if grammarToken == grammar.R_NumericLiteral {
		grammarToken = grammar.R_INTEGER
	}

	return &tokenNumericLiteral{
		Offsets:     r.commitForTextOffsetRange(uncommitted),
		GrammarRule: grammarToken,
		Decoded:     string(uncommitted),
	}, nil
}
