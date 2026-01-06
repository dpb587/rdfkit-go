package fileresource

import (
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type readerParams struct{}

var _ rdfiotypes.Params = &readerParams{}

func (f *readerParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{}
}

func (f *readerParams) ApplyDefaults() {}
