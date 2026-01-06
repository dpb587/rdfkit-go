package encodingdefaults

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingNquads struct{}

var _ encodingref.RegistryEncoding = &encodingNquads{}

type encodingNquadsDecoderFlags struct {
	CaptureTextOffsets *bool
}

func (e encodingNquads) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingNquadsDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := nquads.DecoderConfig{}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	decoder, err := nquads.NewDecoder(wrapReader(rr, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}

type encodingNquadsEncoderFlags struct {
	Ascii *bool
}

func (e encodingNquads) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	flags := encodingNquadsEncoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := nquads.EncoderConfig{}

	if flags.Ascii != nil {
		options = options.SetASCII(*flags.Ascii)
	}

	encoder, err := nquads.NewEncoder(wrapWriter(ww, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.EncoderHandle{
		Writer:  ww,
		Encoder: encoder,
	}, nil
}
