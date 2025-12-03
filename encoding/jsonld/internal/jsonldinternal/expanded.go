package jsonldinternal

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
)

type ExpandedValue interface {
	isExpandedValue()

	// only used for testing
	AsBuiltin() any
}

//

type ExpandedArray struct {
	Values []ExpandedValue
}

func (*ExpandedArray) isExpandedValue() {}

func (e *ExpandedArray) AsBuiltin() any {
	builtinValues := make([]any, 0, len(e.Values))

	for _, v := range e.Values {
		builtinValues = append(builtinValues, v.AsBuiltin())
	}

	return builtinValues
}

//

type ExpandedObject struct {
	Members               map[string]ExpandedValue
	SourceOffsets         *cursorio.TextOffsetRange
	PropertySourceOffsets *cursorio.TextOffsetRange
}

func (*ExpandedObject) isExpandedValue() {}

func (e *ExpandedObject) AsBuiltin() any {
	builtinMembers := make(map[string]any)

	for k, v := range e.Members {
		builtinMembers[k] = v.AsBuiltin()
	}

	return builtinMembers
}

//

type ExpandedScalarPrimitive struct {
	Value                 inspectjson.Value
	PropertySourceOffsets *cursorio.TextOffsetRange
}

func (*ExpandedScalarPrimitive) isExpandedValue() {}

func (e *ExpandedScalarPrimitive) AsBuiltin() any {
	return e.Value.AsBuiltin()
}
