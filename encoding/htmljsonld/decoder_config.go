package htmljsonld

import (
	"github.com/dpb587/inspectjson-go/inspectjson"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

type DecoderConfig struct {
	nestedErrorListener func(err error)
	parserOptions       inspectjson.ParserOptionsApplier
	documentLoader      jsonldtype.DocumentLoader
}

var _ DecoderOption = DecoderConfig{}

func (b DecoderConfig) SetNestedErrorListener(v func(err error)) DecoderConfig {
	b.nestedErrorListener = v

	return b
}

func (b DecoderConfig) SetParserOptions(v inspectjson.ParserOptionsApplier) DecoderConfig {
	b.parserOptions = v

	return b
}

func (b DecoderConfig) SetDocumentLoader(v jsonldtype.DocumentLoader) DecoderConfig {
	b.documentLoader = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.nestedErrorListener != nil {
		s.nestedErrorListener = b.nestedErrorListener
	}

	if b.parserOptions != nil {
		s.parserOptions = b.parserOptions
	}

	if b.documentLoader != nil {
		s.documentLoader = b.documentLoader
	}
}

func (b DecoderConfig) newDecoder(doc *encodinghtml.Document) (*Decoder, error) {
	return &Decoder{
		doc:                 doc,
		docProfile:          doc.GetInfo(),
		nestedErrorListener: b.nestedErrorListener,
		parserOptions:       b.parserOptions,
		documentLoader:      b.documentLoader,
	}, nil
}
