package nquadsrdfio

import (
	"github.com/dpb587/kvstrings-go/kvstrings/kvref"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type encoderParams struct {
	Ascii *bool
}

var _ rdfiotypes.Params = &encoderParams{}

func (f *encoderParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{
		"ascii": kvref.BoolPtr(&f.Ascii, rdfiotypes.ParamMeta{
			Usage: "Use escape sequences for non-ASCII characters",
		}),
	}
}

func (f *encoderParams) ApplyDefaults() {}
