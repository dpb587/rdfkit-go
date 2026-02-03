package jsonld

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/iri/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type EncoderConfig struct {
	base     *string
	prefixes iri.PrefixMappingList
	buffered *bool

	jsonPrefix     *string
	jsonIndent     *string
	jsonEscapeHTML *bool

	bnStringProvider blanknodes.StringProvider
}

func (s EncoderConfig) SetBase(v string) EncoderConfig {
	s.base = &v

	return s
}

func (s EncoderConfig) SetPrefixes(v iri.PrefixMappingList) EncoderConfig {
	s.prefixes = v

	return s
}

func (s EncoderConfig) SetBuffered(v bool) EncoderConfig {
	s.buffered = &v

	return s
}

func (s EncoderConfig) SetBlankNodeStringProvider(v blanknodes.StringProvider) EncoderConfig {
	s.bnStringProvider = v

	return s
}

func (s EncoderConfig) SetIndent(prefix, indent string) EncoderConfig {
	s.jsonPrefix = &prefix
	s.jsonIndent = &indent

	return s
}

func (s EncoderConfig) SetEscapeHTML(v bool) EncoderConfig {
	s.jsonEscapeHTML = &v

	return s
}

func (s EncoderConfig) apply(d *EncoderConfig) {
	if s.base != nil {
		d.base = s.base
	}

	if s.prefixes != nil {
		d.prefixes = s.prefixes
	}

	if s.buffered != nil {
		d.buffered = s.buffered
	}

	if s.bnStringProvider != nil {
		d.bnStringProvider = s.bnStringProvider
	}

	if s.jsonPrefix != nil {
		d.jsonPrefix = s.jsonPrefix
	}

	if s.jsonIndent != nil {
		d.jsonIndent = s.jsonIndent
	}

	if s.jsonEscapeHTML != nil {
		d.jsonEscapeHTML = s.jsonEscapeHTML
	}
}

func (s EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	e := &Encoder{
		w:                json.NewEncoder(w),
		prefixes:         iriutil.NewUsagePrefixMapper(iri.NewPrefixManager(s.prefixes)),
		bnStringProvider: s.bnStringProvider,
		builder:          rdfdescription.NewDatasetResourceListBuilder(),
	}

	if s.base != nil {
		baseIRI, err := iri.ParseBaseIRI(string(*s.base))
		if err != nil {
			return nil, fmt.Errorf("parse base: %v", err)
		}

		e.base = baseIRI
	}

	if s.buffered != nil {
		e.buffered = *s.buffered
	}

	if e.bnStringProvider == nil {
		e.bnStringProvider = blanknodes.NewInt64StringProvider("b%d")
	}

	if s.jsonPrefix != nil && s.jsonIndent != nil {
		e.w.SetIndent(*s.jsonPrefix, *s.jsonIndent)
	}

	if s.jsonEscapeHTML != nil {
		e.w.SetEscapeHTML(*s.jsonEscapeHTML)
	}

	return e, nil
}
