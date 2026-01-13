package rdfiotypes

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type EncoderOptions struct {
	Type    string
	BaseIRI rdf.IRI
	Params  []string
	Patcher GenericOptionsPatcherFunc

	DecoderPipe *DecoderHandle
}

func (base EncoderOptions) ApplyOptions(r Registry, ww Writer, opts *EncoderOptions) error {
	if len(base.Type) > 0 {
		opts.Type = base.Type
	}

	if len(base.BaseIRI) > 0 {
		opts.BaseIRI = base.BaseIRI
	}

	opts.Params = append(opts.Params, base.Params...)

	if base.Patcher != nil {
		// TODO merge?
		opts.Patcher = base.Patcher
	}

	if base.DecoderPipe != nil {
		opts.DecoderPipe = base.DecoderPipe
	}

	return nil
}

type EncoderHandle struct {
	Writer  Writer
	Encoder encoding.Encoder
}

func (b *EncoderHandle) Close() error {
	if err := b.Encoder.Close(); err != nil {
		return err
	}

	if err := b.Writer.Close(); err != nil {
		return err
	}

	return nil
}

func (dh *EncoderHandle) GetQuadsEncoder() encoding.QuadsEncoder {
	switch d := dh.Encoder.(type) {
	case encoding.QuadsEncoder:
		return d
	case encoding.TriplesEncoder:
		return encodingutil.QuadAsTripleEncoder{
			TriplesEncoder: d,
		}
	}

	panic(fmt.Errorf("unexpected encoder type: %T", dh.Encoder))
}

//

func PropagateDecoderPipeBlankNodeStringProvider(h *DecoderHandle) blanknodes.StringProvider {
	if h != nil && h.DecoderBlankNodes != nil {
		if bnStringProvider, ok := h.DecoderBlankNodes.(blanknodes.StringProviderProvider); ok {
			return bnStringProvider.GetStringProvider(blanknodes.NewUUIDStringProvider("", nil))
		}
	}

	return nil
}
