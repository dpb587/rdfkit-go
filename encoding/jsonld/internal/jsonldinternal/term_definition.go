package jsonldinternal

import (
	"net/url"

	"github.com/dpb587/inspectjson-go/inspectjson"
)

type TermDefinition struct {
	// [4.1] an IRI mapping (IRI),
	// [dpb] commonly referred to as "IRI mapping"; TODO rename?
	IRI      ExpandedIRI
	IRIValue inspectjson.Value

	// [4.1] a prefix flag (boolean),
	// [dpb] commonly referred to as "prefix flag"; TODO rename?
	Prefix bool

	// [4.1] a protected flag (boolean),
	Protected bool

	// [4.1] a reverse property flag (boolean),
	ReverseProperty bool

	// [4.1] an optional base URL (IRI),
	BaseURL *url.URL

	// [4.1] an optional context (context),
	// [dpb] sometimes referred to as "local context"?
	Context inspectjson.Value

	// [4.1] an optional container mapping (array of strings),
	ContainerMapping []string

	// [4.1] an optional direction mapping ("ltr" or "rtl"),
	// DirectionMapping      *string
	// string or null
	DirectionMappingValue inspectjson.Value // inspectjson.NullValue | inspectjson.StringValue

	// [4.1] an optional index mapping (string),
	IndexMapping *ExpandedIRIasIRI

	// [4.1] an optional language mapping (string),
	// LanguageMapping      *string
	// string or null
	LanguageMappingValue inspectjson.Value // inspectjson.NullValue | inspectjson.StringValue

	// [4.1] an optional nest value (string),
	NestValue *string

	// [4.1] and an optional type mapping (IRI).
	TypeMapping      ExpandedIRI
	TypeMappingValue inspectjson.Value
}

func (d *TermDefinition) Equals(d2 *TermDefinition) bool {
	if d == nil && d2 == nil {
		// equal
	} else if d != nil && d2 != nil {
		if !d.IRI.Equals(d2.IRI) {
			return false
		}
	} else {
		return false
	}

	if d.Prefix != d2.Prefix {
		return false
	}

	if d.Protected != d2.Protected {
		return false
	}

	if d.ReverseProperty != d2.ReverseProperty {
		return false
	}

	if d.BaseURL == nil && d2.BaseURL == nil {
		// equal
	} else if d.BaseURL != nil && d2.BaseURL != nil {
		if d.BaseURL.String() != d2.BaseURL.String() {
			return false
		}
	} else {
		return false
	}

	// TODO okay to ignore Context?

	if len(d.ContainerMapping) != len(d2.ContainerMapping) {
		return false
	}

	for i, v := range d.ContainerMapping {
		if v != d2.ContainerMapping[i] {
			return false
		}
	}

	if d.DirectionMappingValue == nil && d2.DirectionMappingValue == nil {
		// equal
	} else if d.DirectionMappingValue != nil && d2.DirectionMappingValue != nil {
		if d.DirectionMappingValue.AsBuiltin() != d2.DirectionMappingValue.AsBuiltin() {
			return false
		}
	} else {
		return false
	}

	if d.IndexMapping == nil && d2.IndexMapping == nil {
		// equal
	} else if d.IndexMapping != nil && d2.IndexMapping != nil {
		if !d.IndexMapping.Equals(d2.IndexMapping) {
			return false
		}
	} else {
		return false
	}

	if d.LanguageMappingValue == nil && d2.LanguageMappingValue == nil {
		// equal
	} else if d.LanguageMappingValue != nil && d2.LanguageMappingValue != nil {
		if d.LanguageMappingValue.AsBuiltin() != d2.LanguageMappingValue.AsBuiltin() {
			return false
		}
	} else {
		return false
	}

	if d.NestValue == nil && d2.NestValue == nil {
		// equal
	} else if d.NestValue != nil && d2.NestValue != nil {
		if *d.NestValue != *d2.NestValue {
			return false
		}
	} else {
		return false
	}

	if d.TypeMapping == nil && d2.TypeMapping == nil {
		// equal
	} else if d.TypeMapping != nil && d2.TypeMapping != nil {
		if !d.TypeMapping.Equals(d2.TypeMapping) {
			return false
		}
	} else {
		return false
	}

	return true
}
