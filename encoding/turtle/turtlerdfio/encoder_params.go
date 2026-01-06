package turtlerdfio

import (
	"github.com/dpb587/kvstrings-go/kvstrings/kvref"
	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type encoderParams struct {
	Buffered        *bool
	IrisUseBase     *bool
	IrisUsePrefixes []string
	Resources       *bool
}

var _ rdfiotypes.Params = &encoderParams{}

func (f *encoderParams) NewParamsCollection() rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{
		"buffered": kvref.BoolPtr(&f.Buffered, rdfiotypes.ParamMeta{
			Usage: "Load all statements into memory before writing any output",
		}),
		"iris.useBase": kvref.BoolPtr(&f.IrisUseBase, rdfiotypes.ParamMeta{
			Usage: "Prefer IRIs relative to the resource IRI",
		}),
		"iris.usePrefix": kvref.StringList(&f.IrisUsePrefixes, rdfiotypes.ParamMeta{
			Usage: "Prefer IRIs using a prefix. Use the syntax of \"{prefix}:{iri}\", \"rdfa-context\", or \"none\"",
		}),
		"resources": kvref.BoolPtr(&f.Resources, rdfiotypes.ParamMeta{
			Usage: "Write nested statements and resource descriptions (implies buffered=true)",
		}),
	}
}

func (f *encoderParams) ApplyDefaults() {
	if f.Buffered == nil {
		f.Buffered = ptr.Value(true)
	}

	if f.IrisUseBase == nil {
		f.IrisUseBase = ptr.Value(true)
	}

	if len(f.IrisUsePrefixes) == 0 {
		f.IrisUsePrefixes = []string{"rdfa-context"}
	}

	if f.Resources == nil {
		f.Resources = ptr.Value(false)
	}
}
