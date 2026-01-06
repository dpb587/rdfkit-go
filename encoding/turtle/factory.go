package turtle

import (
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/turtle/turtlecontent"
)

type Factory struct {
	opts FactoryOptions
}

type FactoryOptions struct {
	EncoderOptions []EncoderOption
	DecoderOptions []DecoderOption
}

var _ encoding.TriplesFactory = &Factory{}

func NewFactory(opts FactoryOptions) *Factory {
	return &Factory{
		opts: opts,
	}
}

func (e *Factory) NewEncoder(w io.Writer) (encoding.TriplesEncoder, error) {
	return NewEncoder(w, e.opts.EncoderOptions...)
}

func (e *Factory) NewDecoder(r io.Reader) (encoding.TriplesDecoder, error) {
	return NewDecoder(r, e.opts.DecoderOptions...)
}

func (e *Factory) GetContentMetadata() encoding.ContentMetadata {
	return turtlecontent.DefaultMetadata
}
