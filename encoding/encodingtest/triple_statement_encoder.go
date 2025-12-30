package encodingtest

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type TriplesEncoderOptions struct {
	BlankNodeStringer        blanknodeutil.Stringer
	BlankNodeStringMapperVar string
	Source                   []byte
}

type TriplesEncoder struct {
	w    io.Writer
	opts TriplesEncoderOptions
}

var _ encoding.TriplesEncoder = &TriplesEncoder{}

func NewTriplesEncoder(w io.Writer, opts TriplesEncoderOptions) *TriplesEncoder {
	if opts.BlankNodeStringer == nil {
		opts.BlankNodeStringer = blanknodeutil.NewStringerInt64()
	}

	if opts.BlankNodeStringMapperVar == "" {
		opts.BlankNodeStringMapperVar = "bnTODO"
	}

	ww := &TriplesEncoder{
		w:    w,
		opts: opts,
	}

	w.Write([]byte("encodingtest.TripleStatementList{\n"))

	return ww
}

func (w *TriplesEncoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{
		FileExt:   ".go",
		MediaType: "application/octet-stream",
	}
}

func (w *TriplesEncoder) Close() error {
	w.w.Write([]byte("}\n"))

	return nil
}

func (w *TriplesEncoder) AddTriple(ctx context.Context, t rdf.Triple) error {
	return w.AddTripleStatement(ctx, TripleStatement{
		Triple: t,
	})
}

func (w *TriplesEncoder) AddTripleStatement(_ context.Context, s TripleStatement) error {
	w.w.Write([]byte("\t{\n"))

	{
		w.w.Write([]byte("\t\tTriple: rdf.Triple{\n"))

		{
			w.w.Write([]byte("\t\t\tSubject: "))

			switch subject := s.Triple.Subject.(type) {
			case rdf.BlankNode:
				fmt.Fprintf(w.w, "%s.MapBlankNodeIdentifier(%q)", w.opts.BlankNodeStringMapperVar, w.opts.BlankNodeStringer.GetBlankNodeIdentifier(subject))
			case rdf.IRI:
				fmt.Fprintf(w.w, "rdf.IRI(%q)", subject)
			default:
				return fmt.Errorf("unsupported subject type: %T", subject)
			}

			w.w.Write([]byte(",\n"))
		}

		{
			w.w.Write([]byte("\t\t\tPredicate: "))

			switch predicate := s.Triple.Predicate.(type) {
			case rdf.IRI:
				fmt.Fprintf(w.w, "rdf.IRI(%q)", predicate)
			default:
				return fmt.Errorf("unsupported predicate type: %T", predicate)
			}

			w.w.Write([]byte(",\n"))
		}

		{
			w.w.Write([]byte("\t\t\tObject: "))

			switch object := s.Triple.Object.(type) {
			case rdf.BlankNode:
				fmt.Fprintf(w.w, "%s.MapBlankNodeIdentifier(%q)", w.opts.BlankNodeStringMapperVar, w.opts.BlankNodeStringer.GetBlankNodeIdentifier(object))
			case rdf.IRI:
				fmt.Fprintf(w.w, "rdf.IRI(%q)", object)
			case rdf.Literal:
				fmt.Fprintf(w.w, "rdf.Literal{\n")
				fmt.Fprintf(w.w, "\t\t\t\tLexicalForm: %q,\n", object.LexicalForm)
				fmt.Fprintf(w.w, "\t\t\t\tDatatype: rdf.IRI(%q),\n", object.Datatype)

				if object.Tag != nil {
					switch tag := object.Tag.(type) {
					case rdf.LanguageLiteralTag:
						fmt.Fprintf(w.w, "\t\t\t\tTag: rdf.LanguageLiteralTag{\n")
						fmt.Fprintf(w.w, "\t\t\t\t\tLanguage: %q,\n", tag.Language)
						fmt.Fprintf(w.w, "\t\t\t\t},\n")
					case rdf.DirectionalLanguageLiteralTag:
						fmt.Fprintf(w.w, "\t\t\t\tTag: rdf.DirectionalLanguageLiteralTag{\n")
						fmt.Fprintf(w.w, "\t\t\t\t\tLanguage: %q,\n", tag.Language)
						fmt.Fprintf(w.w, "\t\t\t\t\tBaseDirection: %q,\n", tag.BaseDirection)
						fmt.Fprintf(w.w, "\t\t\t\t},\n")
					}
				}

				w.w.Write([]byte("\t\t\t}"))
			default:
				return fmt.Errorf("unsupported object type: %T", object)
			}

			w.w.Write([]byte(",\n"))
		}

		w.w.Write([]byte("\t\t},\n"))
	}

	if len(s.TextOffsets) > 0 {
		w.w.Write([]byte("\t\tTextOffsets: encoding.StatementTextOffsets{\n"))

		if v, ok := s.TextOffsets[encoding.SubjectStatementOffsets]; ok {
			w.writeOffsetRange("\t\t\t", "Subject", v)
		}

		if v, ok := s.TextOffsets[encoding.PredicateStatementOffsets]; ok {
			w.writeOffsetRange("\t\t\t", "Predicate", v)
		}

		if v, ok := s.TextOffsets[encoding.ObjectStatementOffsets]; ok {
			w.writeOffsetRange("\t\t\t", "Object", v)
		}

		w.w.Write([]byte("\t\t},\n"))
	}

	w.w.Write([]byte("\t},\n"))

	return nil
}

func (w *TriplesEncoder) AddQuad(ctx context.Context, q rdf.Quad) error {
	return w.AddTripleStatement(ctx, TripleStatement{
		Triple: q.Triple,
	})
}

func (w *TriplesEncoder) AddQuadStatement(ctx context.Context, s QuadStatement) error {
	return w.AddTripleStatement(ctx, TripleStatement{
		Triple:      s.Quad.Triple,
		TextOffsets: s.TextOffsets,
	})
}

// var reNL = regexp.MustCompile(`\r?\n`)

func (w *TriplesEncoder) writeOffsetRange(indent, field string, r cursorio.OffsetRange) {
	if len(w.opts.Source) > 0 {
		sourceRaw := reNL.Split(string(w.opts.Source[r.OffsetRangeFrom().ByteOffset():r.OffsetRangeUntil().ByteOffset()]), -1)

		fmt.Fprintf(w.w, "%s// ", indent)

		if len(sourceRaw) > 1 {
			fl, ll := sourceRaw[0], sourceRaw[len(sourceRaw)-1]
			if len(fl) > 32 {
				fl = strings.TrimRightFunc(fl[:32], unicode.IsSpace)
			}

			if len(ll) > 32 {
				ll = strings.TrimLeftFunc(ll[len(ll)-32:], unicode.IsSpace)
			}

			fmt.Fprintf(w.w, "%s ... %s", fl, ll)
		} else {
			al := sourceRaw[0]

			if len(al) > 69 {
				al = strings.TrimRightFunc(al[:32], unicode.IsSpace) + " ... " + strings.TrimLeftFunc(al[len(al)-32:], unicode.IsSpace)
			}

			fmt.Fprintf(w.w, "%s", al)
		}

		w.w.Write([]byte("\n"))
	}

	fmt.Fprintf(w.w, "%sencoding.%sStatementOffsets: ", indent, field)

	switch rr := r.(type) {
	case cursorio.ByteOffsetRange:
		fmt.Fprintf(w.w, "cursorio.ByteOffsetRange{%d, %d}", rr.From, rr.Until)
	case cursorio.TextOffsetRange:
		fmt.Fprintf(w.w, "cursorio.TextOffsetRange{\n")
		fmt.Fprintf(w.w, "%s\tFrom: cursorio.TextOffset{Byte: %d, LineColumn: cursorio.TextLineColumn{%d, %d}},\n", indent, rr.From.Byte, rr.From.LineColumn[0], rr.From.LineColumn[1])
		fmt.Fprintf(w.w, "%s\tUntil: cursorio.TextOffset{Byte: %d, LineColumn: cursorio.TextLineColumn{%d, %d}},\n", indent, rr.Until.Byte, rr.Until.LineColumn[0], rr.Until.LineColumn[1])
		fmt.Fprintf(w.w, "%s}", indent)
	default:
		panic(fmt.Errorf("unsupported offset range type: %T", rr))
	}

	w.w.Write([]byte(",\n"))
}
