package encodingtest

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

const QuadsEncoderContentTypeIdentifier encoding.ContentTypeIdentifier = "internal.dev.quads"

type QuadsEncoderOptions struct {
	BlankNodeStringProvider   blanknodes.StringProvider
	BlankNodeStringFactoryVar string
	Source                    []byte
}

type QuadsEncoder struct {
	w    io.Writer
	opts QuadsEncoderOptions
}

var _ encoding.QuadsEncoder = &QuadsEncoder{}

func NewQuadsEncoder(w io.Writer, opts QuadsEncoderOptions) *QuadsEncoder {
	if opts.BlankNodeStringProvider == nil {
		opts.BlankNodeStringProvider = blanknodes.NewInt64StringProvider("b%d")
	}

	if opts.BlankNodeStringFactoryVar == "" {
		opts.BlankNodeStringFactoryVar = "bnTODO"
	}

	ww := &QuadsEncoder{
		w:    w,
		opts: opts,
	}

	w.Write([]byte("encodingtest.QuadStatementList{\n"))

	return ww
}

func (w *QuadsEncoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return QuadsEncoderContentTypeIdentifier
}

func (w *QuadsEncoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{
		FileExt: ".go",
		MediaType: encoding.ContentMediaType{
			Type:    "application",
			Subtype: "octet-stream",
		},
	}
}

func (w *QuadsEncoder) Close() error {
	w.w.Write([]byte("}\n"))

	return nil
}

func (w *QuadsEncoder) AddTriple(ctx context.Context, t rdf.Triple) error {
	return w.AddQuadStatement(ctx, QuadStatement{
		Quad: rdf.Quad{
			Triple:    t,
			GraphName: nil,
		},
	})
}

func (w *QuadsEncoder) AddTripleStatement(ctx context.Context, s TripleStatement) error {
	return w.AddQuadStatement(ctx, QuadStatement{
		Quad: rdf.Quad{
			Triple:    s.Triple,
			GraphName: nil,
		},
		TextOffsets: s.TextOffsets,
	})
}

func (w *QuadsEncoder) AddQuad(ctx context.Context, t rdf.Quad) error {
	return w.AddQuadStatement(ctx, QuadStatement{
		Quad: t,
	})
}

func (w *QuadsEncoder) AddQuadStatement(_ context.Context, s QuadStatement) error {
	w.w.Write([]byte("\t{\n"))

	{
		w.w.Write([]byte("\t\tQuad: rdf.Quad{\n"))

		{
			w.w.Write([]byte("\t\t\tTriple: rdf.Triple{\n"))

			{
				w.w.Write([]byte("\t\t\t\tSubject: "))

				switch subject := s.Quad.Triple.Subject.(type) {
				case rdf.BlankNode:
					fmt.Fprintf(w.w, "%s.NewStringBlankNode(%q)", w.opts.BlankNodeStringFactoryVar, w.opts.BlankNodeStringProvider.GetBlankNodeString(subject))
				case rdf.IRI:
					fmt.Fprintf(w.w, "rdf.IRI(%q)", subject)
				default:
					return fmt.Errorf("unsupported subject type: %T", subject)
				}

				w.w.Write([]byte(",\n"))
			}

			{
				w.w.Write([]byte("\t\t\t\tPredicate: "))

				switch predicate := s.Quad.Triple.Predicate.(type) {
				case rdf.IRI:
					fmt.Fprintf(w.w, "rdf.IRI(%q)", predicate)
				default:
					return fmt.Errorf("unsupported predicate type: %T", predicate)
				}

				w.w.Write([]byte(",\n"))
			}

			{
				w.w.Write([]byte("\t\t\t\tObject: "))

				switch object := s.Quad.Triple.Object.(type) {
				case rdf.BlankNode:
					fmt.Fprintf(w.w, "%s.NewStringBlankNode(%q)", w.opts.BlankNodeStringFactoryVar, w.opts.BlankNodeStringProvider.GetBlankNodeString(object))
				case rdf.IRI:
					fmt.Fprintf(w.w, "rdf.IRI(%q)", object)
				case rdf.Literal:
					fmt.Fprintf(w.w, "rdf.Literal{\n")
					fmt.Fprintf(w.w, "\t\t\t\t\tLexicalForm: %q,\n", object.LexicalForm)
					fmt.Fprintf(w.w, "\t\t\t\t\tDatatype: rdf.IRI(%q),\n", object.Datatype)

					if object.Tag != nil {
						switch tag := object.Tag.(type) {
						case rdf.LanguageLiteralTag:
							fmt.Fprintf(w.w, "\t\t\t\t\tTag: rdf.LanguageLiteralTag{\n")
							fmt.Fprintf(w.w, "\t\t\t\t\t\tLanguage: %q,\n", tag.Language)
							fmt.Fprintf(w.w, "\t\t\t\t\t},\n")
						case rdf.DirectionalLanguageLiteralTag:
							fmt.Fprintf(w.w, "\t\t\t\t\tTag: rdf.DirectionalLanguageLiteralTag{\n")
							fmt.Fprintf(w.w, "\t\t\t\t\t\tLanguage: %q,\n", tag.Language)
							fmt.Fprintf(w.w, "\t\t\t\t\t\tBaseDirection: %q,\n", tag.BaseDirection)
							fmt.Fprintf(w.w, "\t\t\t\t\t},\n")
						}
					}

					w.w.Write([]byte("\t\t\t\t}"))
				default:
					return fmt.Errorf("unsupported object type: %T", object)
				}

				w.w.Write([]byte(",\n"))
			}

			w.w.Write([]byte("\t\t\t},\n"))
		}

		if s.Quad.GraphName != nil {
			w.w.Write([]byte("\t\t\tGraphName: "))

			switch graphName := s.Quad.GraphName.(type) {
			case rdf.BlankNode:
				fmt.Fprintf(w.w, "%s.NewStringBlankNode(%q)", w.opts.BlankNodeStringFactoryVar, w.opts.BlankNodeStringProvider.GetBlankNodeString(graphName))
			case rdf.IRI:
				fmt.Fprintf(w.w, "rdf.IRI(%q)", graphName)
			default:
				return fmt.Errorf("unsupported graph name type: %T", graphName)
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

		if v, ok := s.TextOffsets[encoding.GraphNameStatementOffsets]; ok {
			w.writeOffsetRange("\t\t\t", "GraphName", v)
		}

		w.w.Write([]byte("\t\t},\n"))
	}

	w.w.Write([]byte("\t},\n"))

	return nil
}

var reNL = regexp.MustCompile(`\r?\n`)

func (w *QuadsEncoder) writeOffsetRange(indent, field string, r cursorio.OffsetRange) {
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
