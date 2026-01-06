package fileresource

import (
	"maps"

	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

type readerParams struct {
	Filter *rdfioutil.FilterReader
}

var _ rdfiotypes.Params = &readerParams{}

func newReaderParams() *readerParams {
	return &readerParams{
		Filter: &rdfioutil.FilterReader{},
	}
}

func (f *readerParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	c := rdfiotypes.ParamsCollection{}
	maps.Copy(c, f.Filter.NewParamsCollection("filter"))

	return c
}

func (f *readerParams) ApplyDefaults() {
	f.Filter.ApplyDefaults()
}
