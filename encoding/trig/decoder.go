package trig

import (
	"errors"
	"io"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdf"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(r io.Reader) (*Decoder, error)
}

type statement struct {
	quad        rdf.Quad
	textOffsets encoding.StatementTextOffsets
}

type Decoder struct {
	buf *cursorioutil.RuneBuffer
	doc *cursorio.TextWriter

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
	buildTextOffsets        encodingutil.TextOffsetsBuilderFunc

	stack []readerStack

	err error

	statements []statement
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
	return trigcontent.TypeIdentifier
}

func (r *Decoder) Close() error {
	return nil // TODO more cleanup?
}

func (r *Decoder) Err() error {
	return r.err
}

func (r *Decoder) Next() bool {
	if len(r.statements) > 0 {
		r.statements = r.statements[1:]
	}

	var rsNext readerStack

	for {
		if r.err != nil {
			return false
		} else if len(r.statements) > 0 {
			if rsNext.fn != nil {
				r.pushState(rsNext.ectx, rsNext.fn)
			}

			return true
		} else if rsNext.fn == nil {
			if len(r.stack) == 0 {
				return false
			}

			rsNext = r.stack[len(r.stack)-1]
			r.stack = r.stack[:len(r.stack)-1]
		}

		rsNext, r.err = r.scan(rsNext.ectx, rsNext.fn)
	}
}

func (r *Decoder) Quad() rdf.Quad {
	return r.statements[0].quad
}

func (r *Decoder) Statement() rdf.Statement {
	return r.Quad()
}

func (r *Decoder) StatementTextOffsets() encoding.StatementTextOffsets {
	return r.statements[0].textOffsets
}

type readerStack struct {
	ectx evaluationContext
	fn   scanFunc
}

type scanFunc func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error)

func (r *Decoder) pushState(ectx evaluationContext, fn scanFunc) error {
	r.stack = append(r.stack, readerStack{ectx: ectx, fn: fn})

	return nil
}

func (r *Decoder) scan(ectx evaluationContext, fn scanFunc) (readerStack, error) {
	var uncommitted cursorio.DecodedRuneList

	for {
		r0, err := r.buf.NextRune()
		if err != nil {
			return fn(r, ectx, r0, err)
		}

		switch r0.Rune {
		case '#':
			uncommitted = append(uncommitted, r0)

			for {
				r1, err := r.buf.NextRune()
				if err != nil {
					if errors.Is(err, io.EOF) {
						r.commit(uncommitted.AsDecodedRunes())

						return r.terminate()
					}

					return readerStack{}, err
				}

				uncommitted = append(uncommitted, r1)

				if r1.Rune == '\n' {
					break
				}
			}
		case 0x20, 0x09, 0x0A, 0x0D:
			uncommitted = append(uncommitted, r0)
		default:
			if unicode.IsSpace(r0.Rune) {
				uncommitted = append(uncommitted, r0)

				continue
			} else {
				r.commit(uncommitted.AsDecodedRunes())
			}

			return fn(r, ectx, r0, err)
		}
	}
}

func (r *Decoder) terminate() (readerStack, error) {
	r.stack = nil

	return readerStack{}, nil
}

func (r *Decoder) emit(bl ...statement) (readerStack, error) {
	r.statements = append(r.statements, bl...)

	return readerStack{}, nil
}
