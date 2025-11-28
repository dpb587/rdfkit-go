package jsonld

import (
	"fmt"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

type DecoderConfig struct {
	defaultBase *string

	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset

	parserOptions []inspectjson.ParserOption

	processingMode *string
	documentLoader jsonldtype.DocumentLoader
	rdfDirection   *string

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
}

func (b DecoderConfig) SetDefaultBase(v string) DecoderConfig {
	b.defaultBase = &v

	return b
}

func (b DecoderConfig) SetCaptureTextOffsets(v bool) DecoderConfig {
	b.captureTextOffsets = &v

	return b
}

func (b DecoderConfig) SetInitialTextOffset(v cursorio.TextOffset) DecoderConfig {
	t := true

	b.captureTextOffsets = &t
	b.initialTextOffset = &v

	return b
}

func (b DecoderConfig) SetParserOptions(v ...inspectjson.ParserOption) DecoderConfig {
	b.parserOptions = v

	return b
}

func (b DecoderConfig) SetProcessingMode(v string) DecoderConfig {
	b.processingMode = &v

	return b
}

func (b DecoderConfig) SetDocumentLoader(v jsonldtype.DocumentLoader) DecoderConfig {
	b.documentLoader = v

	return b
}

func (b DecoderConfig) SetRDFDirection(v string) DecoderConfig {
	b.rdfDirection = &v

	return b
}

func (b DecoderConfig) SetBaseDirectiveListener(v DecoderEvent_BaseDirective_ListenerFunc) DecoderConfig {
	b.baseDirectiveListener = v

	return b
}

func (b DecoderConfig) SetPrefixDirectiveListener(v DecoderEvent_PrefixDirective_ListenerFunc) DecoderConfig {
	b.prefixDirectiveListener = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.defaultBase != nil {
		s.defaultBase = b.defaultBase
	}

	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}

	if b.parserOptions != nil {
		s.parserOptions = b.parserOptions
	}

	if b.processingMode != nil {
		s.processingMode = b.processingMode
	}

	if b.documentLoader != nil {
		s.documentLoader = b.documentLoader
	}

	if b.rdfDirection != nil {
		s.rdfDirection = b.rdfDirection
	}

	if b.baseDirectiveListener != nil {
		s.baseDirectiveListener = b.baseDirectiveListener
	}

	if b.prefixDirectiveListener != nil {
		s.prefixDirectiveListener = b.prefixDirectiveListener
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	d := &Decoder{
		r:                       r,
		statementsIdx:           -1,
		documentLoader:          b.documentLoader,
		parserOptions:           b.parserOptions,
		baseDirectiveListener:   b.baseDirectiveListener,
		prefixDirectiveListener: b.prefixDirectiveListener,
		buildTextOffsets:        encodingutil.BuildTextOffsetsNil,
	}

	if b.defaultBase != nil {
		d.defaultBase = *b.defaultBase
	}

	if b.captureTextOffsets != nil {
		d.captureTextOffsets = *b.captureTextOffsets

		if b.initialTextOffset != nil {
			d.initialTextOffset = *b.initialTextOffset
		}

		d.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	if b.processingMode != nil {
		d.processingMode = *b.processingMode
	}

	if b.rdfDirection != nil {
		switch *b.rdfDirection {
		case "i18n-datatype", "compound-literal":
		// good
		default:
			return nil, fmt.Errorf("rdf direction: invalid value: %v", *b.rdfDirection)
		}

		d.rdfDirection = *b.rdfDirection
	}

	return d, nil
}
