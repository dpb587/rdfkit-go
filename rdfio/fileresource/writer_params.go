package fileresource

import (
	"maps"

	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

type writerParams struct {
	Filter *rdfioutil.FilterWriter
}

var _ rdfiotypes.Params = &writerParams{}

func newWriterParams() *writerParams {
	return &writerParams{
		Filter: &rdfioutil.FilterWriter{},
	}
}

func (f *writerParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	c := rdfiotypes.ParamsCollection{}
	maps.Copy(c, f.Filter.NewParamsCollection("filter"))

	return c
}

func (f *writerParams) ApplyDefaults() {
	f.Filter.ApplyDefaults()
}
