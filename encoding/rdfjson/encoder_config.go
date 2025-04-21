package rdfjson

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type EncoderConfig struct {
	prefix            *string
	indent            *string
	blankNodeStringer blanknodeutil.Stringer
}

func (s EncoderConfig) SetPrefix(v string) EncoderConfig {
	s.prefix = &v

	return s
}

func (s EncoderConfig) SetIndent(v string) EncoderConfig {
	s.indent = &v

	return s
}

func (s EncoderConfig) SetBlankNodeStringer(v blanknodeutil.Stringer) EncoderConfig {
	s.blankNodeStringer = v

	return s
}

func (s EncoderConfig) apply(d *EncoderConfig) {
	if s.prefix != nil {
		d.prefix = s.prefix
	}

	if s.indent != nil {
		d.indent = s.indent
	}

	if s.blankNodeStringer != nil {
		d.blankNodeStringer = s.blankNodeStringer
	}
}

func (s EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	ww := &Encoder{
		w:                 w,
		blankNodeStringer: s.blankNodeStringer,
		buf:               map[string]map[string][]any{},
	}

	if s.prefix != nil {
		ww.prefix = *s.prefix
	}

	if s.indent != nil {
		ww.indent = *s.indent
	}

	if ww.blankNodeStringer == nil {
		ww.blankNodeStringer = blanknodeutil.NewStringerInt64()
	}

	return ww, nil
}
