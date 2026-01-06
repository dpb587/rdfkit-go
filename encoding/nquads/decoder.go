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
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(r io.Reader) (*Decoder, error)
}

type Decoder struct {
	buf              *cursorioutil.RuneBuffer
	doc              *cursorio.TextWriter
	bnStringFactory  blanknodes.StringFactory
	buildTextOffsets encodingutil.TextOffsetsBuilderFunc

	err error

	currentQuad        rdf.Quad
	currentTextOffsets encoding.StatementTextOffsets
}

var _ encoding.QuadsDecoder = &Decoder{}
var _ encoding.StatementTextOffsetsProvider = &Decoder{}

func NewDecoder(r io.Reader, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(r)
}

func (r *Decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return nquadscontent.TypeIdentifier
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
		if r.currentQuad.Triple.Subject != nil {
			for {
				r0, err := r.buf.NextRune()
				if err != nil {
					if errors.Is(err, io.EOF) {
						// TODO technically at least one triple must be present
						r.currentQuad = rdf.Quad{}

						return nil
					}

					return grammar.R_nquadsDoc.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
				}

				switch {
				case r0.Rune == '#':
					err = r.drainLine(cursorio.DecodedRuneList{r0})
					if err != nil {
						if errors.Is(err, io.EOF) {
							// TODO technically at least one triple must be present
							r.currentQuad = rdf.Quad{}

							return nil
						}

						return grammar.R_nquadsDoc.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
					}

					goto QUAD_START
				case r0.Rune == 0xD || r0.Rune == 0xA:
					r.commit(r0.AsDecodedRunes())

					goto QUAD_START
				default:
					if unicode.IsSpace(r0.Rune) {
						r.commit(r0.AsDecodedRunes())
					} else {
						return grammar.R_nquadsDoc.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
					}
				}
			}
		}

	QUAD_START:

		subject, subjectRange, err := r.captureSubjectOrGraphValue(grammar.R_subject)
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.currentQuad = rdf.Quad{}

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
				return grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
			}

			switch {
			case r0.Rune == '.':
				r.commit(r0.AsDecodedRunes())

				goto QUAD_DONE
			case r0.Rune == '#':
				err = r.drainLine(cursorio.DecodedRuneList{r0})
				if err != nil {
					return err
				}
			case unicode.IsSpace(r0.Rune):
				r.commit(r0.AsDecodedRunes())
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
				return grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
			}

			switch {
			case r0.Rune == '.':
				goto QUAD_DONE
			case r0.Rune == '#':
				err = r.drainLine(cursorio.DecodedRuneList{r0})
				if err != nil {
					return err
				}
			case unicode.IsSpace(r0.Rune):
				r.commit(r0.AsDecodedRunes())
			default:
				return grammar.R_statement.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
			}
		}

	QUAD_DONE:

		r.currentQuad = rdf.Quad{
			Triple: rdf.Triple{
				Subject:   subject,
				Predicate: predicate,
				Object:    object,
			},
			GraphName: graphName,
		}

		r.currentTextOffsets = r.buildTextOffsets(
			encoding.GraphNameStatementOffsets, graphNameRange,
			encoding.SubjectStatementOffsets, subjectRange,
			encoding.PredicateStatementOffsets, predicateRange,
			encoding.ObjectStatementOffsets, objectRange,
		)

		return nil
	})()
	if err != nil {
		r.err = err

		return false
	}

	return r.currentQuad.Triple.Subject != nil
}

func (r *Decoder) Quad() rdf.Quad {
	return r.currentQuad
}

func (r *Decoder) Statement() rdf.Statement {
	return r.Quad()
}

func (r *Decoder) StatementTextOffsets() encoding.StatementTextOffsets {
	return r.currentTextOffsets
}
