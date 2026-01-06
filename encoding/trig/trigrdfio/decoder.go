package trigrdfio

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/trig"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type decoder struct{}

var _ rdfiotypes.DecoderManager = decoder{}

func NewDecoder() rdfiotypes.DecoderManager {
	return decoder{}
}

func (decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return trigcontent.TypeIdentifier
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

	options := trig.DecoderConfig{}.
		SetDefaultBase(string(opts.BaseIRI))

	if params.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*params.CaptureTextOffsets)
	}

	handle := &rdfiotypes.DecoderHandle{
		Reader: rr,
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]trig.DecoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	decoder, err := trig.NewDecoder(rr, allOptions...)
	if err != nil {
		return nil, err
	}

	handle.Decoder = decoder

	return handle, nil
}
