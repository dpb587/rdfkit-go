package ntriples

import (
	"bytes"
	"io"

	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type EncoderConfig struct {
	ascii             *bool
	blankNodeStringer blanknodeutil.Stringer
}

func (o EncoderConfig) SetASCII(v bool) EncoderConfig {
	o.ascii = &v

	return o
}

func (o EncoderConfig) SetBlankNodeStringer(v blanknodeutil.Stringer) EncoderConfig {
	o.blankNodeStringer = v

	return o
}

func (o EncoderConfig) apply(d *EncoderConfig) {
	if o.ascii != nil {
		d.ascii = o.ascii
	}

	if o.blankNodeStringer != nil {
		d.blankNodeStringer = o.blankNodeStringer
	}
}

func (o EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	ww := &Encoder{
		w:                 w,
		blankNodeStringer: o.blankNodeStringer,
		buf:               bytes.NewBuffer(make([]byte, 0, 4096)),
	}

	if o.ascii != nil && *o.ascii {
		ww.ascii = *o.ascii
	}

	if ww.blankNodeStringer == nil {
		ww.blankNodeStringer = blanknodeutil.NewStringerInt64()
	}

	return ww, nil
}
