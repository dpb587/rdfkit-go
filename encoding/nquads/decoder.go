package nquads

import (
	"errors"
	"io"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/nquads/internal/grammar"
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

var _ rdfio.DatasetStatementIterator = &Decoder{}

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

					return grammar.R_nquadsDoc.Err(r.newOffsetError(err, nil, nil))
				}

				switch {
				case r0 == '#':
					err = r.drainLine([]rune{r0})
					if err != nil {
						if errors.Is(err, io.EOF) {
							// TODO technically at least one triple must be present
							r.currentStatement = nil

							return nil
						}

						return grammar.R_nquadsDoc.Err(r.newOffsetError(err, nil, nil))
					}

					goto QUAD_START
				case r0 == 0xD || r0 == 0xA:
					r.commit([]rune{r0})

					goto QUAD_START
				default:
					if unicode.IsSpace(r0) {
						r.commit([]rune{r0})
					} else {
						return grammar.R_nquadsDoc.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
					}
				}
			}
		}

	QUAD_START:

		subject, subjectRange, err := r.captureSubjectOrGraphValue(grammar.R_subject)
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.currentStatement = nil

				return nil
			}

			return grammar.R_statement.Err(err)
		}

		predicate, predicateRange, err := r.capturePredicate()
		if err != nil {
			return grammar.R_statement.Err(err)
		}

		object, objectRange, err := r.captureObject()
		if err != nil {
			return grammar.R_statement.Err(err)
		}

		var graphName rdf.GraphNameValue
		var graphNameRange *cursorio.TextOffsetRange

		for {
			r0, err := r.buf.NextRune()
			if err != nil {
				return grammar.R_statement.Err(r.newOffsetError(err, nil, nil))
			}

			switch {
			case r0 == '.':
				r.commit([]rune{r0})

				goto QUAD_DONE
			case r0 == '#':
				err = r.drainLine([]rune{r0})
				if err != nil {
					return err
				}
			case unicode.IsSpace(r0):
				r.commit([]rune{r0})
			default:
				r.buf.BacktrackRunes(r0)

				rawGraphName, rawGraphNameRange, err := r.captureSubjectOrGraphValue(grammar.R_graphLabel)
				if err != nil {
					return grammar.R_statement.Err(err)
				}

				graphName = rawGraphName.(rdf.GraphNameValue)
				graphNameRange = rawGraphNameRange

				goto GRAPH_NAME_DONE
			}
		}

	GRAPH_NAME_DONE:

		for {
			r0, err := r.buf.NextRune()
			if err != nil {
				return grammar.R_statement.Err(r.newOffsetError(err, nil, nil))
			}

			switch {
			case r0 == '.':
				goto QUAD_DONE
			case r0 == '#':
				err = r.drainLine([]rune{r0})
				if err != nil {
					return err
				}
			case unicode.IsSpace(r0):
				r.commit([]rune{r0})
			default:
				return grammar.R_statement.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
			}
		}

	QUAD_DONE:

		r.currentStatement = &statement{
			graphName: graphName,
			triple: rdf.Triple{
				Subject:   subject,
				Predicate: predicate,
				Object:    object,
			},
			offsets: r.buildTextOffsets(
				encoding.GraphNameStatementOffsets, graphNameRange,
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

func (r *Decoder) GetGraphName() rdf.GraphNameValue {
	return r.currentStatement.graphName
}

func (r *Decoder) GetTriple() rdf.Triple {
	return r.currentStatement.triple
}

func (r *Decoder) GetStatement() rdfio.Statement {
	return r.currentStatement
}
