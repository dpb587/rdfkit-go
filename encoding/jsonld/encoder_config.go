package jsonld

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type EncoderConfig struct {
	base     *string
	prefixes iriutil.PrefixMap
	buffered *bool

	jsonPrefix     *string
	jsonIndent     *string
	jsonEscapeHTML *bool

	blankNodeStringer blanknodeutil.Stringer
}

func (s EncoderConfig) SetBase(v string) EncoderConfig {
	s.base = &v

	return s
}

func (s EncoderConfig) SetPrefixes(v iriutil.PrefixMap) EncoderConfig {
	s.prefixes = v

	return s
}

func (s EncoderConfig) SetBuffered(v bool) EncoderConfig {
	s.buffered = &v

	return s
}

func (s EncoderConfig) SetBlankNodeStringer(v blanknodeutil.Stringer) EncoderConfig {
	s.blankNodeStringer = v

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

	if s.blankNodeStringer != nil {
		d.blankNodeStringer = s.blankNodeStringer
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
	var prefixes iriutil.PrefixMap

	if s.prefixes == nil {
		prefixes = iriutil.PrefixMap{}
	} else {
		prefixes = s.prefixes
	}

	e := &Encoder{
		w:                 json.NewEncoder(w),
		prefixes:          iriutil.NewPrefixTracker(prefixes),
		blankNodeStringer: s.blankNodeStringer,
		builder:           rdfdescription.NewDatasetResourceListBuilder(),
	}

	if s.base != nil {
		baseIRI, err := iriutil.ParseBaseIRI(string(*s.base))
		if err != nil {
			return nil, fmt.Errorf("parse base: %v", err)
		}

		e.base = baseIRI
	}

	if s.buffered != nil {
		e.buffered = *s.buffered
	}

	if e.blankNodeStringer == nil {
		e.blankNodeStringer = blanknodeutil.NewStringerInt64()
	}

	if s.jsonPrefix != nil && s.jsonIndent != nil {
		e.w.SetIndent(*s.jsonPrefix, *s.jsonIndent)
	}

	if s.jsonEscapeHTML != nil {
		e.w.SetEscapeHTML(*s.jsonEscapeHTML)
	}

	return e, nil
}
