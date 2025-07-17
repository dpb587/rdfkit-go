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
			return nil, nil, grammar.R_subject.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
		}

		switch {
		case r0.Rune == '<':
			iri, iriRange, err := r.captureOpenIRI(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_subject.Err(err)
			}

			return iri, iriRange, nil
		case r0.Rune == '_':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, nil, grammar.R_subject.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
			}

			if r1.Rune != ':' {
				return nil, nil, grammar.R_subject.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, r0.AsDecodedRunes(), r1.AsDecodedRunes()))
			}

			blankNode, blankNodeRange, err := r.captureOpenBlankNode(cursorio.DecodedRuneList{r0, r1})
			if err != nil {
				return nil, nil, grammar.R_subject.Err(err)
			}

			return blankNode, blankNodeRange, nil
		case r0.Rune == '#':
			err = r.drainLine(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_subject.Err(err)
			}
		case unicode.IsSpace(r0.Rune):
			r.commit(r0.AsDecodedRunes())
		default:
			return nil, nil, grammar.R_subject.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
		}
	}
}

func (r *Decoder) capturePredicate() (rdf.PredicateValue, *cursorio.TextOffsetRange, error) {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_predicate.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
		}

		switch {
		case r0.Rune == '<':
			iri, iriRange, err := r.captureOpenIRI(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_predicate.Err(err)
			}

			return iri, iriRange, nil
		case r0.Rune == '#':
			err = r.drainLine(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_predicate.Err(err)
			}
		case unicode.IsSpace(r0.Rune):
			r.commit(r0.AsDecodedRunes())
		default:
			return nil, nil, grammar.R_predicate.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
		}
	}
}

func (r *Decoder) captureObject() (rdf.ObjectValue, *cursorio.TextOffsetRange, error) {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return nil, nil, grammar.R_object.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
		}

		switch {
		case r0.Rune == '<':
			iri, iriRange, err := r.captureOpenIRI(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}

			return iri, iriRange, nil
		case r0.Rune == '_':
			r1, err := r.buf.NextRune()
			if err != nil {
				return nil, nil, grammar.R_object.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
			}

			if r1.Rune != ':' {
				return nil, nil, grammar.R_object.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, r0.AsDecodedRunes(), r1.AsDecodedRunes()))
			}

			blankNode, blankNodeRange, err := r.captureOpenBlankNode(cursorio.DecodedRuneList{r0, r1})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}

			return blankNode, blankNodeRange, nil
		case r0.Rune == '"':
			literal, literalRange, err := r.captureOpenLiteral(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}

			return literal, literalRange, nil
		case r0.Rune == '#':
			err = r.drainLine(cursorio.DecodedRuneList{r0})
			if err != nil {
				return nil, nil, grammar.R_object.Err(err)
			}
		case unicode.IsSpace(r0.Rune):
			r.commit(r0.AsDecodedRunes())
		default:
			return nil, nil, grammar.R_object.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
		}
	}
}

func (r *Decoder) drainLine(uncommitted cursorio.DecodedRuneList) error {
	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.commit(uncommitted.AsDecodedRunes())
			}

			return err
		}

		if r0.Rune == '\n' {
			r.commit(append(uncommitted, r0).AsDecodedRunes())

			return nil
		}

		uncommitted = append(uncommitted, r0)
	}
}
