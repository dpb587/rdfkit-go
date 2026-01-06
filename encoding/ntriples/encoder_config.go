package ntriples

import (
	"bytes"
	"io"

	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type EncoderConfig struct {
	ascii            *bool
	bnStringProvider blanknodes.StringProvider
}

func (o EncoderConfig) SetASCII(v bool) EncoderConfig {
	o.ascii = &v

	return o
}

func (o EncoderConfig) SetBlankNodeStringProvider(v blanknodes.StringProvider) EncoderConfig {
	o.bnStringProvider = v

	return o
}

func (o EncoderConfig) apply(d *EncoderConfig) {
	if o.ascii != nil {
		d.ascii = o.ascii
	}

	if o.bnStringProvider != nil {
		d.bnStringProvider = o.bnStringProvider
	}
}

func (o EncoderConfig) newEncoder(w io.Writer) (*Encoder, error) {
	ww := &Encoder{
		w:                w,
		bnStringProvider: o.bnStringProvider,
		buf:              bytes.NewBuffer(make([]byte, 0, 4096)),
	}

	if o.ascii != nil && *o.ascii {
		ww.ascii = *o.ascii
	}

	if ww.bnStringProvider == nil {
		ww.bnStringProvider = blanknodes.NewInt64StringProvider("b%d")
	}

	return ww, nil
}
