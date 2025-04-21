package turtle

import (
	"fmt"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type DecoderConfig struct {
	defaultBase     *string
	defaultPrefixes iriutil.PrefixMap

	blankNodeStringMapper blanknodeutil.StringMapper

	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
}

func (b DecoderConfig) SetDefaultBase(v string) DecoderConfig {
	b.defaultBase = &v

	return b
}

func (b DecoderConfig) SetDefaultPrefixes(v iriutil.PrefixMap) DecoderConfig {
	b.defaultPrefixes = v

	return b
}

func (b DecoderConfig) SetBlankNodeStringMapper(v blanknodeutil.StringMapper) DecoderConfig {
	b.blankNodeStringMapper = v

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

	if o.blankNodeStringMapper != nil {
		s.blankNodeStringMapper = o.blankNodeStringMapper
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
	var defaultBase *iriutil.ParsedIRI

	if o.defaultBase != nil {
		var err error

		defaultBase, err = iriutil.ParseIRI(*o.defaultBase)
		if err != nil {
			return nil, fmt.Errorf("base url: %v", err)
		}
	}

	var defaultPrefixes iriutil.PrefixMap

	if o.defaultPrefixes != nil {
		defaultPrefixes = o.defaultPrefixes
	} else {
		defaultPrefixes = iriutil.PrefixMap{}
	}

	var blankNodeStringMapper blanknodeutil.StringMapper

	if o.blankNodeStringMapper != nil {
		blankNodeStringMapper = o.blankNodeStringMapper
	} else {
		blankNodeStringMapper = blanknodeutil.NewStringMapper()
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
						Base:                  defaultBase,
						Prefixes:              defaultPrefixes,
						BlankNodeStringMapper: blankNodeStringMapper,
						BlankNodeFactory:      blanknodeutil.NewFactory(),
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
