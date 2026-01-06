package rdfiotypes

import (
	"errors"
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
)

var ErrUnknownEncoding = errors.New("encoding not supported")

//

//

type DecoderManager interface {
	GetContentTypeIdentifier() encoding.ContentTypeIdentifier
	NewDecoderParams() Params
	NewDecoder(rr Reader, opts DecoderOptions) (*DecoderHandle, error)
}

type EncoderManager interface {
	GetContentTypeIdentifier() encoding.ContentTypeIdentifier
	NewEncoderParams() Params
	NewEncoder(rw Writer, opts EncoderOptions) (*EncoderHandle, error)
}

//

func PatchGenericOptions[T any](in []T, f GenericOptionsPatcherFunc) ([]T, error) {
	if f == nil {
		return in, nil
	}

	out, err := f(in)
	if err != nil {
		return nil, fmt.Errorf("patch: %w", err)
	}

	return out.([]T), nil
}

type GenericOptionsPatcherFunc func(opts any) (any, error)
