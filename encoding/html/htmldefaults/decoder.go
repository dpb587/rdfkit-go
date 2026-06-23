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
	"github.com/dpb587/rdfkit-go/encoding/htmlrdfa"
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
	doc   *html.Document
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

func (d *Decoder) GetDocument() *html.Document {
	return d.doc
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
			d.init()

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

func (d *Decoder) init() {
	err := func() error {
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

		if d.cfg.rootVisitor != nil {
			options = options.SetRootVisitor(d.cfg.rootVisitor)
		}

		htmlDocument, err := html.ParseDocument(d.r, options)
		if err != nil {
			return fmt.Errorf("html: %v", err)
		}

		htmlJsonld, err := htmljsonld.NewDecoder(
			htmlDocument,
			append(
				[]htmljsonld.DecoderOption{
					htmljsonld.DecoderConfig{}.
						SetParserOptions(
							inspectjson.TokenizerConfig{}.
								SetLax(true),
						),
				},
				d.cfg.jsonldOptions...,
			)...,
		)
		if err != nil {
			return fmt.Errorf("htmljsonld: %v", err)
		}

		htmlMicrodata, err := htmlmicrodata.NewDecoder(
			htmlDocument,
			append(
				[]htmlmicrodata.DecoderOption{
					htmlmicrodata.DecoderConfig{}.
						SetVocabularyResolver(htmlmicrodata.ItemtypeVocabularyResolver),
				},
				d.cfg.microdataOptions...,
			)...,
		)
		if err != nil {
			return fmt.Errorf("htmlmicrodata: %v", err)
		}

		htmlRdfa, err := htmlrdfa.NewDecoder(htmlDocument, d.cfg.rdfaOptions...)
		if err != nil {
			return fmt.Errorf("htmlrdfa: %v", err)
		}

		d.doc = htmlDocument
		d.iters = []nestedIterator{
			htmlJsonld,
			encodingutil.NewTripleAsQuadDecoder(htmlMicrodata, nil),
			encodingutil.NewTripleAsQuadDecoder(htmlRdfa, nil),
		}

		return nil
	}()
	if err != nil {
		d.err = err
	}
}
