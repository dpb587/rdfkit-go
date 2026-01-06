package ntriplesrdfio

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/ntriplescontent"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type encoder struct{}

var _ rdfiotypes.EncoderManager = encoder{}

func NewEncoder() rdfiotypes.EncoderManager {
	return encoder{}
}

func (encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return ntriplescontent.TypeIdentifier
}

func (e encoder) NewEncoderParams() rdfiotypes.Params {
	return &encoderParams{}
}

func (e encoder) NewEncoder(ww rdfiotypes.Writer, opts rdfiotypes.EncoderOptions) (*rdfiotypes.EncoderHandle, error) {
	params := &encoderParams{}

	err := rdfiotypes.LoadAndApplyParams(params, opts.Params...)
	if err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	options := ntriples.EncoderConfig{}

	if params.Ascii != nil {
		options = options.SetASCII(*params.Ascii)
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]ntriples.EncoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	encoder, err := ntriples.NewEncoder(ww, allOptions...)
	if err != nil {
		return nil, err
	}

	return &rdfiotypes.EncoderHandle{
		Writer:  ww,
		Encoder: encoder,
	}, nil
}
