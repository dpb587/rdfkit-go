package turtlerdfio

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/encoding/turtle/turtlecontent"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type decoder struct{}

var _ rdfiotypes.DecoderManager = decoder{}

func NewDecoder() rdfiotypes.DecoderManager {
	return decoder{}
}

func (decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return turtlecontent.TypeIdentifier
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

	options := turtle.DecoderConfig{}.
		SetBlankNodeStringFactory(bnFactory).
		SetDefaultBase(string(opts.BaseIRI))

	if params.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*params.CaptureTextOffsets)
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]turtle.DecoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	decoder, err := turtle.NewDecoder(rr, allOptions...)
	if err != nil {
		return nil, err
	}

	return &rdfiotypes.DecoderHandle{
		Reader:            rr,
		Decoder:           decoder,
		DecoderBlankNodes: bnFactory,
	}, nil
}
