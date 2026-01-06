package encodingdefaults

import (
	"fmt"
	"net/http"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingJsonld struct {
	jsonldDocumentLoader jsonldtype.DocumentLoader
}

type encodingJsonldDecoderFlags struct {
	CaptureTextOffsets *bool
	Tokenizer          *encodingJsonldDecoderFlagsTokenizer
}

type encodingJsonldDecoderFlagsTokenizer struct {
	Lax *bool
}

var _ encodingref.RegistryEncoding = &encodingJsonld{}

func (e encodingJsonld) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingJsonldDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := jsonld.DecoderConfig{}.
		SetDocumentLoader(jsonldtype.NewDefaultDocumentLoader(http.DefaultClient))

	if len(opts.IRI) > 0 {
		options = options.SetDefaultBase(string(opts.IRI))
	}

	tokenizerOptions := inspectjson.TokenizerConfig{}

	if flags.Tokenizer != nil {
		if flags.Tokenizer.Lax != nil {
			tokenizerOptions = tokenizerOptions.SetLax(*flags.Tokenizer.Lax)
		}
	}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
		tokenizerOptions = tokenizerOptions.SetSourceOffsets(*flags.CaptureTextOffsets)
	}

	options = options.SetParserOptions(tokenizerOptions)

	decoder, err := jsonld.NewDecoder(wrapReader(rr, opts), options)
	if err != nil {
		return nil, err
	}

	return &encodingref.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}

func (e encodingJsonld) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	return nil, encodingref.ErrEncodingNotSupported
}
