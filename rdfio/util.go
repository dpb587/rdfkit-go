package rdfio

import (
	"github.com/dpb587/rdfkit-go/encoding/trig"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type DirectiveAggregator struct {
	Base           []string
	PrefixMappings iriutil.PrefixMappingList
}

func PatchDirectiveAggregatorOptions(d *DirectiveAggregator, opts any) (any, error) {
	switch opts := opts.(type) {
	case []trig.DecoderOption:
		return append(opts, trig.DecoderConfig{}.
			SetBaseDirectiveListener(func(data trig.DecoderEvent_BaseDirective_Data) {
				d.Base = append(d.Base, data.Value)
			}).
			SetPrefixDirectiveListener(func(data trig.DecoderEvent_PrefixDirective_Data) {
				d.PrefixMappings = append(d.PrefixMappings, iriutil.PrefixMapping{
					Prefix:   data.Prefix,
					Expanded: rdf.IRI(data.Expanded),
				})
			}),
		), nil
	case []turtle.DecoderOption:
		return append(opts, turtle.DecoderConfig{}.
			SetBaseDirectiveListener(func(data turtle.DecoderEvent_BaseDirective_Data) {
				d.Base = append(d.Base, data.Value)
			}).
			SetPrefixDirectiveListener(func(data turtle.DecoderEvent_PrefixDirective_Data) {
				d.PrefixMappings = append(d.PrefixMappings, iriutil.PrefixMapping{
					Prefix:   data.Prefix,
					Expanded: rdf.IRI(data.Expanded),
				})
			}),
		), nil
	}

	return opts, nil
}
