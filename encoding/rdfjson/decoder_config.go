package rdfjson

import (
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
)

type DecoderConfig struct {
	tokenizerOptions []inspectjson.TokenizerOption

	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset
}

func (b DecoderConfig) SetTokenizerOptions(v ...inspectjson.TokenizerOption) DecoderConfig {
	b.tokenizerOptions = v

	return b
}

func (b DecoderConfig) SetCaptureTextOffsets(v bool) DecoderConfig {
	b.captureTextOffsets = &v

	return b
}

func (b DecoderConfig) SetInitialTextOffset(v cursorio.TextOffset) DecoderConfig {
	b.initialTextOffset = &v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.tokenizerOptions != nil {
		s.tokenizerOptions = b.tokenizerOptions
	}

	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	d := &Decoder{
		r:                r,
		statementsIdx:    -1,
		buildTextOffsets: encodingutil.BuildTextOffsetsNil,
	}

	if len(b.tokenizerOptions) > 0 {
		d.topts = append(d.topts, b.tokenizerOptions...)
	}

	if b.captureTextOffsets != nil && *b.captureTextOffsets {
		var initialTextOffset cursorio.TextOffset

		if b.initialTextOffset != nil {
			initialTextOffset = *b.initialTextOffset
		}

		d.topts = append(d.topts, inspectjson.TokenizerConfig{}.SetSourceInitialOffset(initialTextOffset))
		d.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	return d, nil
}
