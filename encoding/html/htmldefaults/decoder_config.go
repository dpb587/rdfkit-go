package htmldefaults

import (
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

type DecoderConfig struct {
	location             *string
	captureTextOffsets   *bool
	initialTextOffset    *cursorio.TextOffset
	jsonldDocumentLoader jsonldtype.DocumentLoader
}

func (b DecoderConfig) SetLocation(v string) DecoderConfig {
	b.location = &v

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

func (b DecoderConfig) SetDocumentLoaderJSONLD(v jsonldtype.DocumentLoader) DecoderConfig {
	b.jsonldDocumentLoader = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.location != nil {
		s.location = b.location
	}

	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}

	if b.jsonldDocumentLoader != nil {
		s.jsonldDocumentLoader = b.jsonldDocumentLoader
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	return &Decoder{
		r:   r,
		cfg: b,
	}, nil
}
