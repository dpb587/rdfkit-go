package rdfiotypes

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf"
)

type DecoderOptions struct {
	Type    string
	BaseIRI rdf.IRI
	Params  []string
	Patcher GenericOptionsPatcherFunc
}

func (next DecoderOptions) ApplyOptions(r Registry, rr Reader, base *DecoderOptions) error {
	if len(next.Type) > 0 {
		base.Type = next.Type
	}

	if len(next.BaseIRI) > 0 {
		base.BaseIRI = next.BaseIRI
	}

	base.Params = append(base.Params, next.Params...)

	if next.Patcher != nil {
		if base.Patcher != nil {
			originalPatcher := base.Patcher

			base.Patcher = GenericOptionsPatcherFunc(func(wopts any) (any, error) {
				nopts, err := originalPatcher(wopts)
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

type DecoderHandle struct {
	Reader            Reader
	Decoder           encoding.Decoder
	DecoderBlankNodes rdf.BlankNodeFactory
}

func (b *DecoderHandle) Close() error {
	if err := b.Decoder.Close(); err != nil {
		return err
	}

	if err := b.Reader.Close(); err != nil {
		return err
	}

	return nil
}

func (dh *DecoderHandle) GetQuadsDecoder() encoding.QuadsDecoder {
	switch d := dh.Decoder.(type) {
	case encoding.QuadsDecoder:
		return d
	case encoding.TriplesDecoder:
		return encodingutil.NewTripleAsQuadDecoder(d, nil)
	}

	panic(fmt.Errorf("unexpected decoder type: %T", dh.Decoder))
}
