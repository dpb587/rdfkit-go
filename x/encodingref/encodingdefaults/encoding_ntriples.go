package encodingdefaults

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingNtriples struct{}

var _ encodingref.RegistryEncoding = &encodingNtriples{}

type encodingNtriplesDecoderFlags struct {
	CaptureTextOffsets *bool
}

func (e encodingNtriples) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingNtriplesDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := ntriples.DecoderConfig{}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	decoder, err := ntriples.NewDecoder(wrapReader(rr, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}

type encodingNtriplesEncoderFlags struct {
	Ascii *bool
}

func (e encodingNtriples) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	flags := encodingNtriplesEncoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := ntriples.EncoderConfig{}

	if flags.Ascii != nil {
		options = options.SetASCII(*flags.Ascii)
	}

	encoder, err := ntriples.NewEncoder(wrapWriter(ww, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.EncoderHandle{
		Writer:  ww,
		Encoder: encoder,
	}, nil
}
