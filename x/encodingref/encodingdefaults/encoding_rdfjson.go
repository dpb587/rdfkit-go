package encodingdefaults

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingRdfjson struct{}

type encodingRdfjsonDecoderFlags struct {
	CaptureTextOffsets *bool
}

var _ encodingref.RegistryEncoding = &encodingRdfjson{}

func (e encodingRdfjson) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingRdfjsonDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := rdfjson.DecoderConfig{}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	decoder, err := rdfjson.NewDecoder(wrapReader(rr, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}

type encodingRdfjsonEncoderFlags struct {
}

func (e encodingRdfjson) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	flags := encodingRdfjsonEncoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := rdfjson.EncoderConfig{}

	encoder, err := rdfjson.NewEncoder(wrapWriter(ww, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.EncoderHandle{
		Writer:  ww,
		Encoder: encoder,
	}, nil
}
