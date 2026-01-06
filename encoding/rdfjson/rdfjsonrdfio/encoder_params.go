package rdfjsonrdfio

import (
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type encoderParams struct{}

var _ rdfiotypes.Params = &encoderParams{}

func (f *encoderParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{}
}

func (f *encoderParams) ApplyDefaults() {}
