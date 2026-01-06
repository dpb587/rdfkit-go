package htmljsonld

import (
	"strings"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/htmljsonld/htmljsonldcontent"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/rdf"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(doc *encodinghtml.Document) (*Decoder, error)
}

type Decoder struct {
	doc        *encodinghtml.Document
	docProfile encodinghtml.DocumentInfo

	nestedErrorListener func(err error)
	parserOptions       []inspectjson.ParserOption
	decoderOptions      []jsonld.DecoderOption

	readers []*jsonld.Decoder

	err error

	currentQuad        rdf.Quad
	currentTextOffsets encoding.StatementTextOffsets
}

var _ encoding.QuadsDecoder = &Decoder{}
var _ encoding.StatementTextOffsetsProvider = &Decoder{}

func NewDecoder(doc *encodinghtml.Document, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(doc)
}

func (r *Decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return htmljsonldcontent.TypeIdentifier
}

func (r *Decoder) Close() error {
	return nil
}

func (w *Decoder) Err() error {
	return w.err
}

func (w *Decoder) Next() bool {
	if w.err != nil {
		return false
	} else if w.readers == nil {
		w.readers = []*jsonld.Decoder{}

		w.walkNode(w.doc.GetRoot())
	}

	for len(w.readers) > 0 {
		if w.readers[0].Next() {
			w.currentQuad = w.readers[0].Quad()
			w.currentTextOffsets = w.readers[0].StatementTextOffsets()

			return true
		} else if err := w.readers[0].Err(); err != nil {
			if w.nestedErrorListener != nil {
				w.nestedErrorListener(err)
			} else {
				w.err = err

				return false
			}
		}

		w.readers = w.readers[1:]
	}

	return false
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

func (w *Decoder) walkNode(n *html.Node) {
	if n.DataAtom == atom.Script {
		for _, attr := range n.Attr {
			if attr != (html.Attribute{Key: "type", Val: "application/ld+json"}) {
				continue
			} else if n.FirstChild == nil {
				// TODO warn?
				return
			}

			dopt := jsonld.DecoderConfig{}.
				SetDefaultBase(w.docProfile.BaseURL).
				SetParserOptions(w.parserOptions...)

			if w.docProfile.HasNodeMetadata {
				if nodeOffsets, ok := w.doc.GetNodeMetadata(n); ok {
					dopt = dopt.
						SetCaptureTextOffsets(true).
						SetInitialTextOffset(nodeOffsets.TokenOffsets.Until)
				}
			}

			nodeReader, err := jsonld.NewDecoder(strings.NewReader(n.FirstChild.Data), append(w.decoderOptions, dopt)...)
			if err != nil {
				if w.nestedErrorListener != nil {
					w.nestedErrorListener(err)
				}

				continue
			}

			w.readers = append(w.readers, nodeReader)

			return
		}

		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w.walkNode(c)
	}
}
