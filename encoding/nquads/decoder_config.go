package nquads

import (
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type DecoderConfig struct {
	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset

	bnStringFactory blanknodes.StringFactory
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

func (b DecoderConfig) SetBlankNodeStringFactory(v blanknodes.StringFactory) DecoderConfig {
	b.bnStringFactory = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}

	if b.bnStringFactory != nil {
		s.bnStringFactory = b.bnStringFactory
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	d := &Decoder{
		buf:              cursorioutil.NewRuneBuffer(r),
		bnStringFactory:  b.bnStringFactory,
		buildTextOffsets: encodingutil.BuildTextOffsetsNil,
	}

	if b.captureTextOffsets != nil && *b.captureTextOffsets {
		var initialTextOffset cursorio.TextOffset

		if b.initialTextOffset != nil {
			initialTextOffset = *b.initialTextOffset
		}

		d.doc = cursorio.NewTextWriter(initialTextOffset)
		d.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	if b.bnStringFactory == nil {
		d.bnStringFactory = blanknodes.NewStringFactory()
	}

	return d, nil
}
