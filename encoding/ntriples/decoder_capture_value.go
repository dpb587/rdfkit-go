package ntriples

import (
	"errors"
	"io"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal/grammar"
	"github.com/dpb587/rdfkit-go/rdf"
)

func (r *Decoder) captureSubject() (rdf.SubjectValue, *cursorio.TextOffsetRange, error) {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_subject.Err(r.newOffsetError(err, nil, nil))
		}

		switch {
		case r0 == '<':
			iri, iriRange, err := r.captureOpenIRI([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_subject.Err(err)
			}

			return iri, iriRange, nil
		case r0 == '_':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, nil, grammar.R_subject.Err(r.newOffsetError(err, nil, nil))
			}

			if r1 != ':' {
				return nil, nil, grammar.R_subject.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, []rune{r0}, []rune{r1}))
			}

			blankNode, blankNodeRange, err := r.captureOpenBlankNode([]rune{r0, r1})
			if err != nil {
				return nil, nil, grammar.R_subject.Err(err)
			}

			return blankNode, blankNodeRange, nil
		case r0 == '#':
			err = r.drainLine([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_subject.Err(err)
			}
		case unicode.IsSpace(r0):
			r.commit([]rune{r0})
		default:
			return nil, nil, grammar.R_subject.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
		}
	}
}

func (r *Decoder) capturePredicate() (rdf.PredicateValue, *cursorio.TextOffsetRange, error) {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_predicate.Err(r.newOffsetError(err, nil, nil))
		}

		switch {
		case r0 == '<':
			iri, iriRange, err := r.captureOpenIRI([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_predicate.Err(err)
			}

			return iri, iriRange, nil
		case r0 == '#':
			err = r.drainLine([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_predicate.Err(err)
			}
		case unicode.IsSpace(r0):
			r.commit([]rune{r0})
		default:
			return nil, nil, grammar.R_predicate.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
		}
	}
}

func (r *Decoder) captureObject() (rdf.ObjectValue, *cursorio.TextOffsetRange, error) {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_object.Err(r.newOffsetError(err, nil, nil))
		}

		switch {
		case r0 == '<':
			iri, iriRange, err := r.captureOpenIRI([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}

			return iri, iriRange, nil
		case r0 == '_':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, nil, grammar.R_object.Err(r.newOffsetError(err, nil, nil))
			}

			if r1 != ':' {
				return nil, nil, grammar.R_object.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, []rune{r0}, []rune{r1}))
			}

			blankNode, blankNodeRange, err := r.captureOpenBlankNode([]rune{r0, r1})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}

			return blankNode, blankNodeRange, nil
		case r0 == '"':
			literal, literalRange, err := r.captureOpenLiteral([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}

			return literal, literalRange, nil
		case r0 == '#':
			err = r.drainLine([]rune{r0})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}
		case unicode.IsSpace(r0):
			r.commit([]rune{r0})
		default:
			return nil, nil, grammar.R_object.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
		}
	}
}

func (r *Decoder) drainLine(uncommitted []rune) error {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.commit(uncommitted)
			}

			return err
		}

		if r0 == '\n' {
			r.commit(append(uncommitted, r0))

			return nil
		}

		uncommitted = append(uncommitted, r0)
	}
}
