package encodingdefaults

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type encodingDev struct {
}

var _ encodingref.RegistryEncoding = &encodingDev{}

func (e encodingDev) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	return nil, encodingref.ErrEncodingNotSupported
}

func (e encodingDev) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	switch cti {
	case encodingtest.DiscardEncoderContentTypeIdentifier:
		return &encodingref.EncoderHandle{
			Writer:  ww,
			Encoder: encodingtest.DiscardEncoder,
		}, nil
	case encodingtest.TriplesEncoderContentTypeIdentifier:
		return &encodingref.EncoderHandle{
			Writer:  ww,
			Encoder: encodingtest.NewTriplesEncoder(wrapWriter(ww, opts), encodingtest.TriplesEncoderOptions{}),
		}, nil
	case encodingtest.QuadsEncoderContentTypeIdentifier:
		return &encodingref.EncoderHandle{
			Writer:  ww,
			Encoder: encodingtest.NewQuadsEncoder(wrapWriter(ww, opts), encodingtest.QuadsEncoderOptions{}),
		}, nil
	}

	return nil, encodingref.ErrEncodingNotSupported
}
