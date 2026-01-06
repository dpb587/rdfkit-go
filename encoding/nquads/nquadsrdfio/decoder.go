package nquadsrdfio

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type decoder struct{}

var _ rdfiotypes.DecoderManager = decoder{}

func NewDecoder() rdfiotypes.DecoderManager {
	return decoder{}
}

func (decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return nquadscontent.TypeIdentifier
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

	options := nquads.DecoderConfig{}

	if params.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*params.CaptureTextOffsets)
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]nquads.DecoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	decoder, err := nquads.NewDecoder(rr, allOptions...)
	if err != nil {
		return nil, err
	}

	return &rdfiotypes.DecoderHandle{
		Reader:  rr,
		Decoder: decoder,
	}, nil
}
