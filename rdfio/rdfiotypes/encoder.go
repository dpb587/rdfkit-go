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

func (next EncoderOptions) ApplyOptions(r Registry, ww Writer, base *EncoderOptions) error {
	if len(next.Type) > 0 {
		base.Type = next.Type
	}

	if len(next.BaseIRI) > 0 {
		base.BaseIRI = next.BaseIRI
	}

	base.Params = append(base.Params, next.Params...)

	if next.Patcher != nil {
		if base.Patcher != nil {
			base.Patcher = GenericOptionsPatcherFunc(func(wopts any) (any, error) {
				nopts, err := base.Patcher(wopts)
				if err != nil {
					return nil, fmt.Errorf("base patcher: %v", err)
				}

				return next.Patcher(nopts)
			})
		} else {
			base.Patcher = next.Patcher
		}
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
