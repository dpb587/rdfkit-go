package jsonld

import "github.com/dpb587/cursorio-go/cursorio"

type DecoderEvent_BaseDirective_ListenerFunc func(data DecoderEvent_BaseDirective_Data)

type DecoderEvent_BaseDirective_Data struct {
	Value      string
	ValueRange *cursorio.TextOffsetRange
}

//

type DecoderEvent_PrefixDirective_ListenerFunc func(data DecoderEvent_PrefixDirective_Data)

type DecoderEvent_PrefixDirective_Data struct {
	Prefix      string
	PrefixRange *cursorio.TextOffsetRange

	Expanded      string
	ExpandedRange *cursorio.TextOffsetRange
}
