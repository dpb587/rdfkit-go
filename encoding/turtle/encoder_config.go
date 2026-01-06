package turtle

import (
	"fmt"
	"io"
	"slices"

	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type EncoderConfig struct {
	base     *string
	prefixes iriutil.PrefixMap

	bnStringProvider blanknodes.StringProvider

	buffered     *bool
	bufferedSort *bool
}

func (s EncoderConfig) SetBase(v string) EncoderConfig {
	s.base = &v

	return s
}

func (s EncoderConfig) SetPrefixes(v iriutil.PrefixMap) EncoderConfig {
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
}

func (s EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	var prefixes iriutil.PrefixMap

	if s.prefixes == nil {
		prefixes = iriutil.PrefixMap{}
	} else {
		prefixes = s.prefixes
	}

	e := &Encoder{
		w:                w,
		prefixes:         iriutil.NewPrefixTracker(prefixes),
		bnStringProvider: s.bnStringProvider,
	}

	if s.base != nil {
		baseIRI, err := iriutil.ParseBaseIRI(string(*s.base))
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

	if !e.buffered && (e.base != nil || len(prefixes) > 0) {
		prefixMappings := prefixes.AsPrefixMappingList()
		slices.SortFunc(prefixMappings, iriutil.ComparePrefixMappingByPrefix)

		var baseString string

		if e.base != nil {
			baseString = e.base.String()
		}

		written, err := WriteDocumentHeader(e.w, baseString, prefixMappings)
		if err != nil {
			return nil, fmt.Errorf("write header: %v", err)
		} else if written > 0 {
			_, err = e.w.Write([]byte("\n"))
			if err != nil {
				return nil, fmt.Errorf("write header: %v", err)
			}
		}
	}

	return e, nil
}
