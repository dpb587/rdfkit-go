package htmljsonld

import (
	"github.com/dpb587/inspectjson-go/inspectjson"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
)

type DecoderConfig struct {
	nestedErrorListener func(err error)
	parserOptions       []inspectjson.ParserOption
	decoderOptions      []jsonld.DecoderOption
}

var _ DecoderOption = DecoderConfig{}

func (b DecoderConfig) SetNestedErrorListener(v func(err error)) DecoderConfig {
	b.nestedErrorListener = v

	return b
}

func (b DecoderConfig) SetParserOptions(v ...inspectjson.ParserOption) DecoderConfig {
	b.parserOptions = v

	return b
}

func (b DecoderConfig) SetDecoderOptions(v ...jsonld.DecoderOption) DecoderConfig {
	b.decoderOptions = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.nestedErrorListener != nil {
		s.nestedErrorListener = b.nestedErrorListener
	}

	if b.parserOptions != nil {
		s.parserOptions = b.parserOptions
	}

	if b.decoderOptions != nil {
		s.decoderOptions = b.decoderOptions
	}
}

func (b DecoderConfig) newDecoder(doc *encodinghtml.Document) (*Decoder, error) {
	return &Decoder{
		doc:                 doc,
		docProfile:          doc.GetInfo(),
		nestedErrorListener: b.nestedErrorListener,
		parserOptions:       b.parserOptions,
		decoderOptions:      b.decoderOptions,
	}, nil
}
