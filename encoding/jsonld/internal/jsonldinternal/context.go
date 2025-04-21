package jsonldinternal

import (
	"net/url"

	"github.com/dpb587/inspectjson-go/inspectjson"
)

type Context struct {
	// [4.1] the active term definitions which specify how keys and values have to be interpreted (array of term definitions),
	TermDefinitions map[string]*TermDefinition

	// [4.1] the current base IRI (IRI),
	BaseURL      *url.URL // TODO BaseURL as rdf.Term and BaseURLURL
	BaseURLValue inspectjson.Value

	// [4.1] the original base URL (IRI),
	OriginalBaseURL *url.URL

	// [4.1] an inverse context (inverse context),
	InverseContext *Context

	// [4.1] an optional vocabulary mapping (IRI),
	VocabularyMapping      ExpandedIRI
	VocabularyMappingValue inspectjson.Value

	// [4.1] an optional default language (string),
	// DefaultLanguage      *string
	DefaultLanguageValue *inspectjson.StringValue

	// [4.1] an optional default base direction ("ltr" or "rtl"),
	// [dpb] commonly referred to as "base direction"; TODO rename?
	// DefaultDirection      *string
	DefaultDirectionValue *inspectjson.StringValue

	// [4.1] and an optional previous context (context), used when a non-propagated context is defined.
	PreviousContext *Context

	_processor *contextProcessor
}

func (c *Context) clone() *Context {
	cClone := &Context{
		TermDefinitions:   map[string]*TermDefinition{},
		BaseURL:           c.BaseURL,
		BaseURLValue:      c.BaseURLValue,
		OriginalBaseURL:   c.OriginalBaseURL,
		InverseContext:    c.InverseContext,
		VocabularyMapping: c.VocabularyMapping,
		// DefaultLanguage:       c.DefaultLanguage,
		DefaultLanguageValue: c.DefaultLanguageValue,
		// DefaultDirection:      c.DefaultDirection,
		DefaultDirectionValue: c.DefaultDirectionValue,
		PreviousContext:       c.PreviousContext,

		_processor: c._processor,
	}

	for k, v := range c.TermDefinitions {
		cClone.TermDefinitions[k] = v
	}

	return cClone
}
