package htmldefaults

import (
	"fmt"
	"io"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/html/htmlcontent"
	"github.com/dpb587/rdfkit-go/encoding/htmljsonld"
	"github.com/dpb587/rdfkit-go/encoding/htmlmicrodata"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/rdfa"
	"github.com/dpb587/rdfkit-go/rdf"
)

type nestedIterator interface {
	encoding.QuadsDecoder
	encoding.StatementTextOffsetsProvider
}

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(r io.Reader) (*Decoder, error)
}

type Decoder struct {
	r     io.Reader
	cfg   DecoderConfig
	err   error
	iters []nestedIterator
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

func (d *Decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return htmlcontent.TypeIdentifier
}

func (d *Decoder) Close() error {
	for _, it := range d.iters {
		it.Close()
	}

	return nil
}

func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) Next() bool {
	for {
		if d.err != nil {
			return false
		} else if d.iters == nil {
			d.iters, d.err = d.init()

			continue
		} else if len(d.iters) == 0 {
			return false
		} else if d.iters[0].Next() {
			return true
		} else if d.iters[0].Err() != nil {
			d.err = fmt.Errorf("decode[%s]: %w", d.iters[0].GetContentTypeIdentifier(), d.iters[0].Err())

			return false
		}

		d.iters = d.iters[1:]
	}
}

func (d *Decoder) Quad() rdf.Quad {
	return d.iters[0].Quad()
}

func (d *Decoder) Statement() rdf.Statement {
	return d.iters[0].Statement()
}

func (d *Decoder) StatementTextOffsets() encoding.StatementTextOffsets {
	return d.iters[0].StatementTextOffsets()
}

func (d *Decoder) init() ([]nestedIterator, error) {
	options := html.DocumentConfig{}

	if d.cfg.location != nil {
		options = options.SetLocation(*d.cfg.location)
	}

	if d.cfg.captureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*d.cfg.captureTextOffsets)
	}

	if d.cfg.initialTextOffset != nil {
		options = options.SetInitialTextOffset(*d.cfg.initialTextOffset)
	}

	htmlDocument, err := html.ParseDocument(d.r, options)
	if err != nil {
		return nil, fmt.Errorf("html: %v", err)
	}

	htmlJsonld, err := htmljsonld.NewDecoder(
		htmlDocument,
		htmljsonld.DecoderConfig{}.
			SetParserOptions(
				inspectjson.TokenizerConfig{}.
					SetLax(true),
			).
			SetDecoderOptions(jsonld.DecoderConfig{}.
				SetDocumentLoader(d.cfg.jsonldDocumentLoader),
			),
	)
	if err != nil {
		return nil, fmt.Errorf("htmljsonld: %v", err)
	}

	htmlMicrodata, err := htmlmicrodata.NewDecoder(
		htmlDocument,
		htmlmicrodata.DecoderConfig{}.
			SetVocabularyResolver(htmlmicrodata.ItemtypeVocabularyResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("htmlmicrodata: %v", err)
	}

	htmlRdfa, err := rdfa.NewDecoder(htmlDocument)
	if err != nil {
		return nil, fmt.Errorf("rdfa: %v", err)
	}

	return []nestedIterator{
		htmlJsonld,
		encodingutil.NewTripleAsQuadDecoder(htmlMicrodata, nil),
		encodingutil.NewTripleAsQuadDecoder(htmlRdfa, nil),
	}, nil
}
