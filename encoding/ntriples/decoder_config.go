package ntriples

import (
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type DecoderConfig struct {
	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset

	blankNodeStringMapper blanknodeutil.StringMapper
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

func (b DecoderConfig) SetBlankNodeStringMapper(v blanknodeutil.StringMapper) DecoderConfig {
	b.blankNodeStringMapper = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}

	if b.blankNodeStringMapper != nil {
		s.blankNodeStringMapper = b.blankNodeStringMapper
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	d := &Decoder{
		buf:                   cursorioutil.NewRuneBuffer(r),
		blankNodeStringMapper: b.blankNodeStringMapper,
		buildTextOffsets:      encodingutil.BuildTextOffsetsNil,
	}

	if b.captureTextOffsets != nil && *b.captureTextOffsets {
		var initialTextOffset cursorio.TextOffset

		if b.initialTextOffset != nil {
			initialTextOffset = *b.initialTextOffset
		}

		d.doc = cursorio.NewTextWriter(initialTextOffset)
		d.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	if d.blankNodeStringMapper == nil {
		d.blankNodeStringMapper = blanknodeutil.NewStringMapper()
	}

	return d, nil
}
