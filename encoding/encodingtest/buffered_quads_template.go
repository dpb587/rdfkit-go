package encodingtest

import (
	"context"
	"fmt"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/terms"
)

type templateExecutor interface {
	Execute(w io.Writer, data any) error
}

type BufferedQuadsTemplateOptions struct {
	Source                      io.Reader
	SourceType                  encoding.ContentTypeIdentifier
	Output                      io.Writer
	OutputContentTypeIdentifier encoding.ContentTypeIdentifier
	OutputContentMetadata       encoding.ContentMetadata
	OutputTemplate              templateExecutor
	Formatter                   terms.Formatter
}

type BufferedQuadsTemplate struct {
	source                      io.Reader
	sourceType                  encoding.ContentTypeIdentifier
	output                      io.Writer
	outputContentTypeIdentifier encoding.ContentTypeIdentifier
	outputContentMetadata       encoding.ContentMetadata
	outputTemplate              templateExecutor
	formatter                   terms.Formatter
	records                     []bufferedQuadsTemplateRecord
}

type bufferedQuadsTemplateRecord struct {
	Encoded     bufferedQuadsTemplateRecordEncoded
	TextOffsets map[string]cursorio.TextOffsetRange
}

type bufferedQuadsTemplateRecordEncoded struct {
	Subject   string
	Predicate string
	Object    string
	GraphName *string
}

var _ encoding.QuadsEncoder = &BufferedQuadsTemplate{}

func NewBufferedQuadsTemplate(opts BufferedQuadsTemplateOptions) *BufferedQuadsTemplate {
	return &BufferedQuadsTemplate{
		source:                      opts.Source,
		sourceType:                  opts.SourceType,
		output:                      opts.Output,
		outputContentTypeIdentifier: opts.OutputContentTypeIdentifier,
		outputContentMetadata:       opts.OutputContentMetadata,
		outputTemplate:              opts.OutputTemplate,
		formatter:                   opts.Formatter,
	}
}

func (w *BufferedQuadsTemplate) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return w.outputContentTypeIdentifier
}

func (w *BufferedQuadsTemplate) GetContentMetadata() encoding.ContentMetadata {
	return w.outputContentMetadata
}

func (w *BufferedQuadsTemplate) Close() error {
	data := struct {
		Source     []byte
		SourceType encoding.ContentTypeIdentifier
		Quads      []bufferedQuadsTemplateRecord
	}{
		SourceType: w.sourceType,
		Quads:      w.records,
	}

	if w.source != nil {
		b, err := io.ReadAll(w.source)
		if err != nil {
			return fmt.Errorf("read input: %v", err)
		}

		data.Source = b
	}

	return w.outputTemplate.Execute(w.output, data)
}

func (w *BufferedQuadsTemplate) AddTriple(ctx context.Context, t rdf.Triple) error {
	return w.AddQuadStatement(ctx, QuadStatement{
		Quad: rdf.Quad{
			Triple:    t,
			GraphName: nil,
		},
	})
}

func (w *BufferedQuadsTemplate) AddTripleStatement(ctx context.Context, s TripleStatement) error {
	return w.AddQuadStatement(ctx, QuadStatement{
		Quad: rdf.Quad{
			Triple:    s.Triple,
			GraphName: nil,
		},
		TextOffsets: s.TextOffsets,
	})
}

func (w *BufferedQuadsTemplate) AddQuad(ctx context.Context, t rdf.Quad) error {
	return w.AddQuadStatement(ctx, QuadStatement{
		Quad: t,
	})
}

func (w *BufferedQuadsTemplate) AddQuadStatement(_ context.Context, s QuadStatement) error {
	encoded := bufferedQuadsTemplateRecord{
		Encoded: bufferedQuadsTemplateRecordEncoded{
			Subject:   w.formatter.FormatTerm(s.Quad.Triple.Subject),
			Predicate: w.formatter.FormatTerm(s.Quad.Triple.Predicate),
			Object:    w.formatter.FormatTerm(s.Quad.Triple.Object),
		},
	}

	if s.Quad.GraphName != nil {
		encoded.Encoded.GraphName = ptr.Value(w.formatter.FormatTerm(s.Quad.GraphName))
	}

	if s.TextOffsets != nil {
		encoded.TextOffsets = map[string]cursorio.TextOffsetRange{}

		if v, ok := s.TextOffsets[encoding.SubjectStatementOffsets]; ok {
			encoded.TextOffsets["Subject"] = v
		}

		if v, ok := s.TextOffsets[encoding.PredicateStatementOffsets]; ok {
			encoded.TextOffsets["Predicate"] = v
		}

		if v, ok := s.TextOffsets[encoding.ObjectStatementOffsets]; ok {
			encoded.TextOffsets["Object"] = v
		}

		if v, ok := s.TextOffsets[encoding.GraphNameStatementOffsets]; ok {
			encoded.TextOffsets["GraphName"] = v
		}
	}

	w.records = append(w.records, encoded)

	return nil
}
