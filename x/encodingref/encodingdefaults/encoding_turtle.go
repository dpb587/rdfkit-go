package encodingdefaults

import (
	"context"
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionutil"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingTurtle struct{}

type encodingTurtleDecoderFlags struct {
	CaptureTextOffsets *bool
}

var _ encodingref.RegistryEncoding = &encodingTurtle{}

func (e encodingTurtle) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingTurtleDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := turtle.DecoderConfig{}

	if len(opts.IRI) > 0 {
		options = options.SetDefaultBase(string(opts.IRI))
	}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	handle := &encodingref.DecoderHandle{
		Reader: rr,
	}

	decoder, err := turtle.NewDecoder(
		wrapReader(rr, opts),
		options,
		turtle.DecoderConfig{}.
			SetBaseDirectiveListener(func(data turtle.DecoderEvent_BaseDirective_Data) {
				handle.DecodedBase = append(handle.DecodedBase, data.Value)
			}).
			SetPrefixDirectiveListener(func(data turtle.DecoderEvent_PrefixDirective_Data) {
				handle.DecodedPrefixMappings = append(handle.DecodedPrefixMappings, iriutil.PrefixMapping{
					Prefix:   data.Prefix,
					Expanded: rdf.IRI(data.Expanded),
				})
			}),
	)
	if err != nil {
		return nil, err
	}

	handle.Decoder = decoder

	return handle, nil
}

type encodingTurtleEncoderFlags struct {
	Buffer          *bool
	DefaultPrefixes *string
	Resources       *bool
}

func (e encodingTurtle) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	flags := encodingTurtleEncoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := turtle.EncoderConfig{}

	if len(opts.IRI) > 0 {
		options = options.SetBase(string(opts.IRI))
	}

	if flags.Buffer == nil || *flags.Buffer {
		options = options.SetBuffered(true)
	}

	if flags.DefaultPrefixes == nil || *flags.DefaultPrefixes == "rdfa-context" {
		options = options.SetPrefixes(iriutil.NewPrefixMap(rdfacontext.InitialContext()...))
	}

	encoder, err := turtle.NewEncoder(wrapWriter(ww, opts), options)
	if err != nil {
		return nil, err
	}

	var wrappedEncoder encoding.Encoder = encoder

	if flags.Resources == nil || *flags.Resources {
		wrappedEncoder = rdfdescriptionutil.NewBufferedTriplesEncoder(
			context.Background(),
			encoder,
			true,
		)
	}

	return &encodingref.EncoderHandle{
		Writer:  ww,
		Encoder: wrappedEncoder,
	}, nil
}
