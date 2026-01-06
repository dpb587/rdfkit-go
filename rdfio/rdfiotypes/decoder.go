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

func (base DecoderOptions) ApplyOptions(r Registry, rr Reader, opts *DecoderOptions) error {
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
