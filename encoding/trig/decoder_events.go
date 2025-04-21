package trig

import "github.com/dpb587/cursorio-go/cursorio"

type DecoderEvent_BaseDirective_ListenerFunc func(data DecoderEvent_BaseDirective_Data)

type DecoderEvent_BaseDirective_Data struct {
	Value        string
	ValueOffsets *cursorio.TextOffsetRange
}

//

type DecoderEvent_PrefixDirective_ListenerFunc func(data DecoderEvent_PrefixDirective_Data)

type DecoderEvent_PrefixDirective_Data struct {
	Prefix        string
	PrefixOffsets *cursorio.TextOffsetRange

	Expanded        string
	ExpandedOffsets *cursorio.TextOffsetRange
}
