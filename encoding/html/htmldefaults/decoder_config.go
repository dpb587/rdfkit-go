package htmldefaults

import (
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/htmljsonld"
	"github.com/dpb587/rdfkit-go/encoding/htmlmicrodata"
	"github.com/dpb587/rdfkit-go/encoding/htmlrdfa"
)

type DecoderConfig struct {
	location           *string
	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset
	jsonldOptions      []htmljsonld.DecoderOption
	microdataOptions   []htmlmicrodata.DecoderOption
	rdfaOptions        []htmlrdfa.DecoderOption
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

func (b DecoderConfig) SetJSONLDOptions(v ...htmljsonld.DecoderOption) DecoderConfig {
	b.jsonldOptions = v

	return b
}

func (b DecoderConfig) SetMicrodataOptions(v ...htmlmicrodata.DecoderOption) DecoderConfig {
	b.microdataOptions = v

	return b
}

func (b DecoderConfig) SetRDFaOptions(v ...htmlrdfa.DecoderOption) DecoderConfig {
	b.rdfaOptions = v

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

	if b.jsonldOptions != nil {
		s.jsonldOptions = b.jsonldOptions
	}

	if b.microdataOptions != nil {
		s.microdataOptions = b.microdataOptions
	}

	if b.rdfaOptions != nil {
		s.rdfaOptions = b.rdfaOptions
	}
}

func (b DecoderConfig) newDecoder(r io.Reader) (*Decoder, error) {
	return &Decoder{
		r:   r,
		cfg: b,
	}, nil
}
