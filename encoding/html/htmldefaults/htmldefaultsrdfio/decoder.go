package htmldefaultsrdfio

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/html/htmlcontent"
	"github.com/dpb587/rdfkit-go/encoding/html/htmldefaults"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
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
	return htmlcontent.TypeIdentifier
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

	options := htmldefaults.DecoderConfig{}.
		SetDocumentLoaderJSONLD(e.jsonldDocumentLoader).
		SetLocation(string(opts.BaseIRI))

	if params.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*params.CaptureTextOffsets)
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]htmldefaults.DecoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	decoder, err := htmldefaults.NewDecoder(rr, allOptions...)
	if err != nil {
		return nil, fmt.Errorf("creating decoder: %v", err)
	}

	return &rdfiotypes.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}
