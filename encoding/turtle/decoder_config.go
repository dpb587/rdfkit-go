package turtle

import (
	"fmt"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type DecoderConfig struct {
	defaultBase     *string
	defaultPrefixes iri.PrefixMappingList

	bnStringFactory blanknodes.StringFactory

	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
}

func (b DecoderConfig) SetDefaultBase(v string) DecoderConfig {
	b.defaultBase = &v

	return b
}

func (b DecoderConfig) SetDefaultPrefixes(v iri.PrefixMappingList) DecoderConfig {
	b.defaultPrefixes = v

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

func (o DecoderConfig) apply(s *DecoderConfig) {
	if o.defaultBase != nil {
		s.defaultBase = o.defaultBase
	}

	if o.defaultPrefixes != nil {
		s.defaultPrefixes = o.defaultPrefixes
	}

	if o.bnStringFactory != nil {
		s.bnStringFactory = o.bnStringFactory
	}

	if o.captureTextOffsets != nil {
		s.captureTextOffsets = o.captureTextOffsets
	}

	if o.initialTextOffset != nil {
		s.initialTextOffset = o.initialTextOffset
	}

	if o.baseDirectiveListener != nil {
		s.baseDirectiveListener = o.baseDirectiveListener
	}

	if o.prefixDirectiveListener != nil {
		s.prefixDirectiveListener = o.prefixDirectiveListener
	}
}

func (o DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	var defaultBase *iri.ParsedIRI

	if o.defaultBase != nil {
		var err error

		defaultBase, err = iri.ParseIRI(*o.defaultBase)
		if err != nil {
			return nil, fmt.Errorf("base url: %v", err)
		}
	}

	var bnStringFactory = o.bnStringFactory

	if o.bnStringFactory == nil {
		bnStringFactory = blanknodes.NewStringFactory()
	}

	d := &Decoder{
		buf:                     cursorioutil.NewRuneBuffer(r),
		baseDirectiveListener:   o.baseDirectiveListener,
		prefixDirectiveListener: o.prefixDirectiveListener,
		buildTextOffsets:        encodingutil.BuildTextOffsetsNil,
		stack: []readerStack{
			{
				ectx: evaluationContext{
					Global: &globalEvaluationContext{
						Base:                   defaultBase,
						Prefixes:               iri.NewPrefixManager(o.defaultPrefixes),
						BlankNodeStringFactory: bnStringFactory,
					},
				},
				fn: reader_scanStatement,
			},
		},
	}

	if o.captureTextOffsets != nil && *o.captureTextOffsets {
		var initialTextOffset cursorio.TextOffset

		if o.initialTextOffset != nil {
			initialTextOffset = *o.initialTextOffset
		}

		d.doc = cursorio.NewTextWriter(initialTextOffset)
		d.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	return d, nil
}
