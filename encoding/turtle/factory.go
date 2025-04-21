package turtle

import (
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
)

type Factory struct {
	opts FactoryOptions
}

type FactoryOptions struct {
	EncoderOptions []EncoderOption
	DecoderOptions []DecoderOption
}

var _ encoding.GraphFactory = &Factory{}

func NewFactory(opts FactoryOptions) *Factory {
	return &Factory{
		opts: opts,
	}
}

func (e *Factory) NewGraphEncoder(w io.Writer) (encoding.GraphEncoder, error) {
	return NewEncoder(w, e.opts.EncoderOptions...)
}

func (e *Factory) NewGraphDecoder(r io.Reader) (encoding.GraphDecoder, error) {
	return NewDecoder(r, e.opts.DecoderOptions...)
}

func (e *Factory) GetGraphEncoderContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{
		FileExt:   ".ttl",
		MediaType: "text/turtle",
		Charset:   "utf-8",
	}
}
