package encodingdefaults

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/trig"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingTrig struct{}

type encodingTrigDecoderFlags struct {
	CaptureTextOffsets *bool
}

var _ encodingref.RegistryEncoding = &encodingTrig{}

func (e encodingTrig) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	flags := encodingTrigDecoderFlags{}

	err := encodingref.UnmarshalFlags(&flags, opts.Flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %v", err)
	}

	options := trig.DecoderConfig{}

	if len(opts.IRI) > 0 {
		options = options.SetDefaultBase(string(opts.IRI))
	}

	if flags.CaptureTextOffsets != nil {
		options = options.SetCaptureTextOffsets(*flags.CaptureTextOffsets)
	}

	handle := &encodingref.DecoderHandle{
		Reader: rr,
	}

	decoder, err := trig.NewDecoder(
		wrapReader(rr, opts),
		options,
		trig.DecoderConfig{}.
			SetBaseDirectiveListener(func(data trig.DecoderEvent_BaseDirective_Data) {
				handle.DecodedBase = append(handle.DecodedBase, data.Value)
			}).
			SetPrefixDirectiveListener(func(data trig.DecoderEvent_PrefixDirective_Data) {
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

func (e encodingTrig) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	return nil, encodingref.ErrEncodingNotSupported
}
