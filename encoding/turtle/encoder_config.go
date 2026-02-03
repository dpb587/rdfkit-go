package turtle

import (
	"fmt"
	"io"
	"slices"

	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/iri/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type EncoderConfig struct {
	base     *string
	prefixes iri.PrefixMappingList

	bnStringProvider blanknodes.StringProvider

	buffered     *bool
	bufferedSort *bool

	baseDirectiveMode   *DirectiveMode
	prefixDirectiveMode *DirectiveMode
}

func (s EncoderConfig) SetBase(v string) EncoderConfig {
	s.base = &v

	return s
}

func (s EncoderConfig) SetPrefixes(v iri.PrefixMappingList) EncoderConfig {
	s.prefixes = v

	return s
}

func (s EncoderConfig) SetBlankNodeStringProvider(v blanknodes.StringProvider) EncoderConfig {
	s.bnStringProvider = v

	return s
}

func (s EncoderConfig) SetBuffered(v bool) EncoderConfig {
	s.buffered = &v

	return s
}

func (s EncoderConfig) SetBufferedSort(v bool) EncoderConfig {
	s.bufferedSort = &v

	return s
}

func (s EncoderConfig) SetBaseDirectiveMode(v DirectiveMode) EncoderConfig {
	s.baseDirectiveMode = &v

	return s
}

func (s EncoderConfig) SetPrefixDirectiveMode(v DirectiveMode) EncoderConfig {
	s.prefixDirectiveMode = &v

	return s
}

func (s EncoderConfig) SetDirectiveMode(v DirectiveMode) EncoderConfig {
	s.baseDirectiveMode = &v
	s.prefixDirectiveMode = &v

	return s
}

func (s EncoderConfig) apply(d *EncoderConfig) {
	if s.base != nil {
		d.base = s.base
	}

	if s.prefixes != nil {
		d.prefixes = s.prefixes
	}

	if s.bnStringProvider != nil {
		d.bnStringProvider = s.bnStringProvider
	}

	if s.buffered != nil {
		d.buffered = s.buffered
	}

	if s.bufferedSort != nil {
		d.bufferedSort = s.bufferedSort
	}

	if s.baseDirectiveMode != nil {
		d.baseDirectiveMode = s.baseDirectiveMode
	}

	if s.prefixDirectiveMode != nil {
		d.prefixDirectiveMode = s.prefixDirectiveMode
	}
}

func (s EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	prefixManager := iri.NewPrefixManager(s.prefixes)

	e := &Encoder{
		w:                   w,
		prefixes:            iriutil.NewUsagePrefixMapper(prefixManager),
		bnStringProvider:    s.bnStringProvider,
		baseDirectiveMode:   DirectiveMode_At,
		prefixDirectiveMode: DirectiveMode_At,
	}

	if s.base != nil {
		baseIRI, err := iri.ParseBaseIRI(string(*s.base))
		if err != nil {
			return nil, fmt.Errorf("parse base: %v", err)
		}

		e.base = baseIRI
	}

	if s.buffered != nil && *s.buffered {
		e.buffered = *s.buffered
		e.bufferedSort = e.buffered
	}

	if s.bufferedSort != nil {
		e.bufferedSort = *s.bufferedSort
	}

	if e.bnStringProvider == nil {
		e.bnStringProvider = blanknodes.NewInt64StringProvider("b%d")
	}

	if !e.buffered && (e.base != nil || len(s.prefixes) > 0) {
		prefixMappings := prefixManager.GetPrefixMappings()
		slices.SortFunc(prefixMappings, iri.ComparePrefixMappingByPrefix)

		var baseString string

		if e.base != nil {
			baseString = e.base.String()
		}

		written, err := WriteDirectives(e.w, WriteDirectivesOptions{
			Base:       baseString,
			Prefixes:   prefixMappings,
			BaseMode:   e.baseDirectiveMode,
			PrefixMode: e.prefixDirectiveMode,
		})
		if err != nil {
			return nil, fmt.Errorf("write header: %v", err)
		} else if written > 0 {
			_, err = e.w.Write([]byte("\n"))
			if err != nil {
				return nil, fmt.Errorf("write header: %v", err)
			}
		}
	}

	if s.baseDirectiveMode != nil {
		e.baseDirectiveMode = *s.baseDirectiveMode
	}

	if s.prefixDirectiveMode != nil {
		e.prefixDirectiveMode = *s.prefixDirectiveMode
	}

	return e, nil
}
