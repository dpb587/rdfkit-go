package rdfiotypes

import "github.com/dpb587/kvstrings-go/kvstrings"

type Params interface {
	NewParamsCollection() ParamsCollection
	ApplyDefaults()
}

type ParamsCollection = kvstrings.Collection[ParamMeta]

type ParamMeta struct {
	Usage      string
	Hidden     bool
	ValueEnums []string
}

func LoadAndApplyParams(p Params, rawParams ...string) error {
	collection := p.NewParamsCollection()

	err := collection.ImportStrings(rawParams...)
	if err != nil {
		return err
	}

	p.ApplyDefaults()

	return nil
}
