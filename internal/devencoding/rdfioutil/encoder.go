package rdfioutil

import (
	"context"
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

type EncoderOptions struct {
	BlankNodeStringer        blanknodeutil.Stringer
	BlankNodeStringMapperVar string
	Source                   []byte
}

type Encoder struct {
	w    io.Writer
	opts EncoderOptions
}

var _ encoding.GraphEncoder = &Encoder{}

func NewEncoder(w io.Writer, opts EncoderOptions) *Encoder {
	if opts.BlankNodeStringer == nil {
		opts.BlankNodeStringer = blanknodeutil.NewStringerInt64()
	}

	if opts.BlankNodeStringMapperVar == "" {
		opts.BlankNodeStringMapperVar = "bnTODO"
	}

	ww := &Encoder{
		w:    w,
		opts: opts,
	}

	w.Write([]byte("rdfio.StatementList{\n"))

	return ww
}

func (w *Encoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{
		FileExt:   ".go",
		MediaType: "application/octet-stream",
	}
}

func (w *Encoder) Close() error {
	w.w.Write([]byte("}\n"))

	return nil
}

func (w *Encoder) PutTriple(ctx context.Context, t rdf.Triple) error {
	return w.PutStatement(ctx, rdfioutil.Statement{
		Triple: t,
	})
}

func (w *Encoder) PutStatement(_ context.Context, s rdfio.Statement) error {
	sGraphName := s.GetGraphName()
	sTriple := s.GetTriple()

	w.w.Write([]byte("\trdfioutil.Statement{\n"))

	if sGraphName != rdf.DefaultGraph {
		w.w.Write([]byte("\t\tGraphName: "))

		switch graphName := sGraphName.(type) {
		case rdf.BlankNode:
			fmt.Fprintf(w.w, "%s.MapBlankNodeIdentifier(%q)", w.opts.BlankNodeStringMapperVar, w.opts.BlankNodeStringer.GetBlankNodeIdentifier(graphName))
		case rdf.IRI:
			fmt.Fprintf(w.w, "rdf.IRI(%q)", graphName)
		default:
			return fmt.Errorf("unsupported graph name type: %T", graphName)
		}

		w.w.Write([]byte(",\n"))
	}

	{
		w.w.Write([]byte("\t\tTriple: rdf.Triple{\n"))

		{
			w.w.Write([]byte("\t\t\tSubject: "))

			switch subject := sTriple.Subject.(type) {
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

			switch predicate := sTriple.Predicate.(type) {
			case rdf.IRI:
				fmt.Fprintf(w.w, "rdf.IRI(%q)", predicate)
			default:
				return fmt.Errorf("unsupported predicate type: %T", predicate)
			}

			w.w.Write([]byte(",\n"))
		}

		{
			w.w.Write([]byte("\t\t\tObject: "))

			switch object := sTriple.Object.(type) {
			case rdf.BlankNode:
				fmt.Fprintf(w.w, "%s.MapBlankNodeIdentifier(%q)", w.opts.BlankNodeStringMapperVar, w.opts.BlankNodeStringer.GetBlankNodeIdentifier(object))
			case rdf.IRI:
				fmt.Fprintf(w.w, "rdf.IRI(%q)", object)
			case rdf.Literal:
				switch object.Datatype {
				case xsdiri.String_Datatype:
					fmt.Fprintf(w.w, "xsdliteral.NewString(%q)", object.LexicalForm)
				case xsdiri.Boolean_Datatype:
					fmt.Fprintf(w.w, "xsdliteral.NewBoolean(%t)", object.LexicalForm == "true")
				default:
					fmt.Fprintf(w.w, "rdf.Literal{\n")
					fmt.Fprintf(w.w, "\t\t\t\tLexicalForm: %q,\n", object.LexicalForm)
					fmt.Fprintf(w.w, "\t\t\t\tDatatype: rdf.IRI(%q),\n", object.Datatype)

					if len(object.Tags) > 0 {
						fmt.Fprintf(w.w, "\t\t\t\tQualifiers: map[rdf.LiteralTag]string{\n")

						tagKeys := slices.Collect(maps.Keys(object.Tags))
						slices.SortFunc(tagKeys, func(i, j rdf.LiteralTag) int {
							if i < j {
								return -1
							}

							return 1
						})

						for _, k := range tagKeys {
							var kString string

							switch k {
							case rdf.BaseDirectionLiteralTag:
								kString = "rdf.BaseDirectionLiteralTag"
							case rdf.LanguageLiteralTag:
								kString = "rdf.LanguageLiteralTag"
							default:
								kString = fmt.Sprintf("%v", k)
							}

							fmt.Fprintf(w.w, "\t\t\t\t\t%s: %q,\n", kString, object.Tags[k])
						}

						fmt.Fprintf(w.w, "\t\t\t\t},\n")
					}

					w.w.Write([]byte("\t\t\t}"))
				}
			default:
				return fmt.Errorf("unsupported object type: %T", object)
			}

			w.w.Write([]byte(",\n"))
		}

		w.w.Write([]byte("\t\t},\n"))
	}

	if sourceRangesProvider, ok := s.(encoding.DecoderTextOffsetsStatement); ok {
		if offsets := sourceRangesProvider.GetDecoderTextOffsets(); len(offsets) > 0 {
			w.w.Write([]byte("\t\tTextOffsets: encoding.StatementTextOffsets{\n"))

			if v, ok := offsets[encoding.GraphNameStatementOffsets]; ok {
				w.writeOffsetRange("\t\t\t", "GraphName", v)
			}

			if v, ok := offsets[encoding.SubjectStatementOffsets]; ok {
				w.writeOffsetRange("\t\t\t", "Subject", v)
			}

			if v, ok := offsets[encoding.PredicateStatementOffsets]; ok {
				w.writeOffsetRange("\t\t\t", "Predicate", v)
			}

			if v, ok := offsets[encoding.ObjectStatementOffsets]; ok {
				w.writeOffsetRange("\t\t\t", "Object", v)
			}

			w.w.Write([]byte("\t\t},\n"))
		}
	}

	w.w.Write([]byte("\t},\n"))

	return nil
}

var reNL = regexp.MustCompile(`\r?\n`)

func (w *Encoder) writeOffsetRange(indent, field string, r cursorio.OffsetRange) {
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
