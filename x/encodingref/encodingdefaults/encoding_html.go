package encodingdefaults

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/html/htmldefaults"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingHtml struct {
	jsonldDocumentLoader jsonldtype.DocumentLoader
}

type encodingHtmlDecoderFlags struct {
	CaptureTextOffsets *bool
}

var _ encodingref.RegistryEncoding = &encodingHtml{}

func (e encodingHtml) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingHtmlDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := htmldefaults.DecoderConfig{}.
		SetDocumentLoaderJSONLD(e.jsonldDocumentLoader)

	if len(opts.IRI) > 0 {
		options = options.SetLocation(string(opts.IRI))
	}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	decoder, err := htmldefaults.NewDecoder(rr, options)
	if err != nil {
		return nil, fmt.Errorf("creating decoder: %v", err)
	}

	return &encodingref.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}

func (e encodingHtml) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	return nil, encodingref.ErrEncodingNotSupported
}
