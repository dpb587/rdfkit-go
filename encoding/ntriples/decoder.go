package ntriples

import (
	"errors"
	"io"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal/grammar"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(r io.Reader) (*Decoder, error)
}

type Decoder struct {
	buf                   *cursorioutil.RuneBuffer
	doc                   *cursorio.TextWriter
	blankNodeStringMapper blanknodeutil.StringMapper
	buildTextOffsets      encodingutil.TextOffsetsBuilderFunc

	err              error
	currentStatement *statement
}

var _ rdfio.GraphStatementIterator = &Decoder{}

func NewDecoder(r io.Reader, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(r)
}

func (r *Decoder) Close() error {
	return nil
}

func (r *Decoder) Err() error {
	return r.err
}

func (r *Decoder) Next() bool {
	if r.err != nil {
		return false
	}

	err := (func() error {
		if r.currentStatement != nil {
			for {
				r0, err := r.buf.NextRune()
				if err != nil {
					if errors.Is(err, io.EOF) {
						// TODO technically at least one triple must be present
						r.currentStatement = nil

						return nil
					}

					return grammar.R_ntriplesDoc.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
				}

				switch {
				case r0.Rune == '#':
					err = r.drainLine(cursorio.DecodedRuneList{r0})
					if err != nil {
						if errors.Is(err, io.EOF) {
							// TODO technically at least one triple must be present
							r.currentStatement = nil

							return nil
						}

						return grammar.R_ntriplesDoc.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
					}

					goto TRIPLE_START
				case r0.Rune == 0xD || r0.Rune == 0xA:
					r.commit(r0.AsDecodedRunes())

					goto TRIPLE_START
				default:
					if unicode.IsSpace(r0.Rune) {
						r.commit(r0.AsDecodedRunes())
					} else {
						return grammar.R_ntriplesDoc.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
					}
				}
			}
		}

	TRIPLE_START:

		subject, subjectRange, err := r.captureSubject()
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.currentStatement = nil

				return nil
			}

			return grammar.R_triple.Err(err)
		}

		predicate, predicateRange, err := r.capturePredicate()
		if err != nil {
			return grammar.R_triple.Err(err)
		}

		object, objectRange, err := r.captureObject()
		if err != nil {
			return grammar.R_triple.Err(err)
		}

		for {
			r0, err := r.buf.NextRune()
			if err != nil {
				return grammar.R_triple.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
			}

			switch {
			case r0.Rune == '.':
				r.commit(r0.AsDecodedRunes())

				goto TRIPLE_DONE
			case r0.Rune == '#':
				err = r.drainLine(cursorio.DecodedRuneList{r0})
				if err != nil {
					return err
				}
			case unicode.IsSpace(r0.Rune):
				r.commit(r0.AsDecodedRunes())
			default:
				return grammar.R_triple.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
			}
		}

	TRIPLE_DONE:

		r.currentStatement = &statement{
			triple: rdf.Triple{
				Subject:   subject,
				Predicate: predicate,
				Object:    object,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, subjectRange,
				encoding.PredicateStatementOffsets, predicateRange,
				encoding.ObjectStatementOffsets, objectRange,
			),
		}

		return nil
	})()
	if err != nil {
		r.err = err

		return false
	}

	return r.currentStatement != nil
}

func (r *Decoder) GetTriple() rdf.Triple {
	return r.currentStatement.triple
}

func (r *Decoder) GetStatement() rdfio.Statement {
	return r.currentStatement
}
