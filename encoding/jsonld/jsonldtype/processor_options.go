package jsonldtype

import "github.com/dpb587/inspectjson-go/inspectjson"

type ProcessorOptions struct {
	ProcessingMode string
	BaseURL        string
	DocumentLoader DocumentLoader
	ExpandContext  inspectjson.Value
}
