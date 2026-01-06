package jsonldrdfio

import (
	"github.com/dpb587/kvstrings-go/kvstrings/kvref"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type decoderParams struct {
	CaptureTextOffsets *bool
	TokenizerLax       *bool
}

var _ rdfiotypes.Params = &decoderParams{}

func (f *decoderParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{
		"captureTextOffsets": kvref.BoolPtr(&f.CaptureTextOffsets, rdfiotypes.ParamMeta{
			Usage: "Capture the line+column offsets for statement properties",
		}),
		"tokenizer.lax": kvref.BoolPtr(&f.TokenizerLax, rdfiotypes.ParamMeta{
			Usage: "Accept and recover common syntax errors",
		}),
	}
}

func (f *decoderParams) ApplyDefaults() {}
