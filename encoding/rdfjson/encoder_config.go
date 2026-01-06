package rdfjson

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type EncoderConfig struct {
	prefix           *string
	indent           *string
	bnStringProvider blanknodes.StringProvider
}

func (s EncoderConfig) SetPrefix(v string) EncoderConfig {
	s.prefix = &v

	return s
}

func (s EncoderConfig) SetIndent(v string) EncoderConfig {
	s.indent = &v

	return s
}

func (s EncoderConfig) SetBlankNodeStringProvider(v blanknodes.StringProvider) EncoderConfig {
	s.bnStringProvider = v

	return s
}

func (s EncoderConfig) apply(d *EncoderConfig) {
	if s.prefix != nil {
		d.prefix = s.prefix
	}

	if s.indent != nil {
		d.indent = s.indent
	}

	if s.bnStringProvider != nil {
		d.bnStringProvider = s.bnStringProvider
	}
}

func (s EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	ww := &Encoder{
		w:                w,
		bnStringProvider: s.bnStringProvider,
		buf:              map[string]map[string][]any{},
	}

	if s.prefix != nil {
		ww.prefix = *s.prefix
	}

	if s.indent != nil {
		ww.indent = *s.indent
	}

	if ww.bnStringProvider == nil {
		ww.bnStringProvider = blanknodes.NewInt64StringProvider("b%d")
	}

	return ww, nil
}
