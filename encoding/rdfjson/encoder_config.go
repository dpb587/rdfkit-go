package rdfjson

import (
	"encoding/json"
	"io"

	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type EncoderConfig struct {
	jsonPrefix     *string
	jsonIndent     *string
	jsonEscapeHTML *bool

	bnStringProvider blanknodes.StringProvider
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

func (s EncoderConfig) SetBlankNodeStringProvider(v blanknodes.StringProvider) EncoderConfig {
	s.bnStringProvider = v

	return s
}

func (s EncoderConfig) apply(d *EncoderConfig) {
	if s.jsonPrefix != nil {
		d.jsonPrefix = s.jsonPrefix
	}

	if s.jsonIndent != nil {
		d.jsonIndent = s.jsonIndent
	}

	if s.jsonEscapeHTML != nil {
		d.jsonEscapeHTML = s.jsonEscapeHTML
	}

	if s.bnStringProvider != nil {
		d.bnStringProvider = s.bnStringProvider
	}
}

func (s EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	ww := &Encoder{
		w:                json.NewEncoder(w),
		bnStringProvider: s.bnStringProvider,
		buf:              map[string]map[string][]any{},
	}

	if ww.bnStringProvider == nil {
		ww.bnStringProvider = blanknodes.NewInt64StringProvider("b%d")
	}

	if s.jsonPrefix != nil && s.jsonIndent != nil {
		ww.w.SetIndent(*s.jsonPrefix, *s.jsonIndent)
	}

	if s.jsonEscapeHTML != nil {
		ww.w.SetEscapeHTML(*s.jsonEscapeHTML)
	}

	return ww, nil
}
