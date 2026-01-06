package encodingdefaults

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingRdfxml struct{}

type encodingRdfxmlDecoderFlags struct {
	CaptureTextOffsets *bool
}

var _ encodingref.RegistryEncoding = &encodingRdfxml{}

func (e encodingRdfxml) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingRdfxmlDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := rdfxml.DecoderConfig{}

	if len(opts.IRI) > 0 {
		options = options.SetBaseURL(string(opts.IRI))
	}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	decoder, err := rdfxml.NewDecoder(wrapReader(rr, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}

func (e encodingRdfxml) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	return nil, encodingref.ErrEncodingNotSupported
}
