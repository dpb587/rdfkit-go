package jsonldrdfio

import (
	"fmt"
	"net/http"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldcontent"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type decoder struct {
	jsonldDocumentLoader jsonldtype.DocumentLoader
}

var _ rdfiotypes.DecoderManager = decoder{}

func NewDecoder(jsonldDocumentLoader jsonldtype.DocumentLoader) rdfiotypes.DecoderManager {
	return decoder{
		jsonldDocumentLoader: jsonldDocumentLoader,
	}
}

func (decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return jsonldcontent.TypeIdentifier
}

func (e decoder) NewDecoderParams() rdfiotypes.Params {
	return &decoderParams{}
}

func (e decoder) NewDecoder(rr rdfiotypes.Reader, opts rdfiotypes.DecoderOptions) (*rdfiotypes.DecoderHandle, error) {
	params := &decoderParams{}

	err := rdfiotypes.LoadAndApplyParams(params, opts.Params...)
	if err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	bnFactory := blanknodes.NewStringFactory()

	options := jsonld.DecoderConfig{}.
		SetBlankNodeStringFactory(bnFactory).
		SetDocumentLoader(jsonldtype.NewDefaultDocumentLoader(http.DefaultClient)).
		SetDefaultBase(string(opts.BaseIRI))

	tokenizerOptions := inspectjson.TokenizerConfig{}

	if params.TokenizerLax != nil {
		tokenizerOptions = tokenizerOptions.SetLax(*params.TokenizerLax)
	}

	if params.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*params.CaptureTextOffsets)
		tokenizerOptions = tokenizerOptions.SetSourceOffsets(*params.CaptureTextOffsets)
	}

	options = options.SetParserOptions(tokenizerOptions)

	allOptions, err := rdfiotypes.PatchGenericOptions([]jsonld.DecoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	decoder, err := jsonld.NewDecoder(rr, allOptions...)
	if err != nil {
		return nil, err
	}

	return &rdfiotypes.DecoderHandle{
		Reader:            rr,
		Decoder:           decoder,
		DecoderBlankNodes: bnFactory,
	}, nil
}
