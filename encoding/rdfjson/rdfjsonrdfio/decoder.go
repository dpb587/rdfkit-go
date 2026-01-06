package rdfjsonrdfio

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson/rdfjsoncontent"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type decoder struct{}

var _ rdfiotypes.DecoderManager = decoder{}

func NewDecoder() rdfiotypes.DecoderManager {
	return decoder{}
}

func (decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return rdfjsoncontent.TypeIdentifier
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

	options := rdfjson.DecoderConfig{}.
		SetBlankNodeStringFactory(bnFactory)

	if params.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*params.CaptureTextOffsets)
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]rdfjson.DecoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	decoder, err := rdfjson.NewDecoder(rr, allOptions...)
	if err != nil {
		return nil, err
	}

	return &rdfiotypes.DecoderHandle{
		Reader:            rr,
		Decoder:           decoder,
		DecoderBlankNodes: bnFactory,
	}, nil
}
