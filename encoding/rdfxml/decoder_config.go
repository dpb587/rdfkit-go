package rdfxml

import (
	"fmt"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type DecoderConfig struct {
	defaultBase     *string
	bnStringFactory blanknodes.StringFactory

	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
	warningListener         func(err error)
}

func (b DecoderConfig) SetDefaultBase(v string) DecoderConfig {
	b.defaultBase = &v

	return b
}

func (b DecoderConfig) SetBlankNodeStringFactory(v blanknodes.StringFactory) DecoderConfig {
	b.bnStringFactory = v

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

func (b DecoderConfig) SetBaseDirectiveListener(v DecoderEvent_BaseDirective_ListenerFunc) DecoderConfig {
	b.baseDirectiveListener = v

	return b
}

func (b DecoderConfig) SetPrefixDirectiveListener(v DecoderEvent_PrefixDirective_ListenerFunc) DecoderConfig {
	b.prefixDirectiveListener = v

	return b
}

func (b DecoderConfig) SetWarningListener(v func(err error)) DecoderConfig {
	b.warningListener = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.defaultBase != nil {
		s.defaultBase = b.defaultBase
	}

	if b.bnStringFactory != nil {
		s.bnStringFactory = b.bnStringFactory
	}

	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}

	if b.baseDirectiveListener != nil {
		s.baseDirectiveListener = b.baseDirectiveListener
	}

	if b.prefixDirectiveListener != nil {
		s.prefixDirectiveListener = b.prefixDirectiveListener
	}

	if b.warningListener != nil {
		s.warningListener = b.warningListener
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	d := &Decoder{
		r:                       r,
		statementsIdx:           -1,
		bnStringFactory:         b.bnStringFactory,
		baseDirectiveListener:   b.baseDirectiveListener,
		prefixDirectiveListener: b.prefixDirectiveListener,
		warningListener:         b.warningListener,
		buildTextOffsets:        encodingutil.BuildTextOffsetsNil,
	}

	if b.defaultBase != nil {
		baseURL, err := iri.ParseIRI(*b.defaultBase)
		if err != nil {
			return nil, fmt.Errorf("base url: %v", err)
		}

		d.baseURL = baseURL
	}

	if b.captureTextOffsets != nil && *b.captureTextOffsets {
		d.captureTextOffsets = true

		if b.initialTextOffset != nil {
			d.initialTextOffset = *b.initialTextOffset
		}

		d.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	if d.bnStringFactory == nil {
		d.bnStringFactory = blanknodes.NewStringFactory()
	}

	return d, nil
}
