package turtle

import (
	"fmt"
	"io"
	"slices"

	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type EncoderConfig struct {
	base     *string
	prefixes iriutil.PrefixMap

	blankNodeStringer blanknodeutil.Stringer

	// Buffered causes all output to be buffered until Close is called. Once closed, the base and any used prefixes will
	// be written, the buffered sections will be sorted and written, and the encoder will no longer be usable.
	buffered *bool
}

func (s EncoderConfig) SetBase(v string) EncoderConfig {
	s.base = &v

	return s
}

func (s EncoderConfig) SetPrefixes(v iriutil.PrefixMap) EncoderConfig {
	s.prefixes = v

	return s
}

func (s EncoderConfig) SetBlankNodeStringer(v blanknodeutil.Stringer) EncoderConfig {
	s.blankNodeStringer = v

	return s
}

func (s EncoderConfig) SetBuffered(v bool) EncoderConfig {
	s.buffered = &v

	return s
}

func (s EncoderConfig) apply(d *EncoderConfig) {
	if s.base != nil {
		d.base = s.base
	}

	if s.prefixes != nil {
		d.prefixes = s.prefixes
	}

	if s.blankNodeStringer != nil {
		d.blankNodeStringer = s.blankNodeStringer
	}

	if s.buffered != nil {
		d.buffered = s.buffered
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
		w:                 w,
		prefixes:          iriutil.NewPrefixTracker(prefixes),
		blankNodeStringer: s.blankNodeStringer,
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
	}

	if e.blankNodeStringer == nil {
		e.blankNodeStringer = blanknodeutil.NewStringerInt64()
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
