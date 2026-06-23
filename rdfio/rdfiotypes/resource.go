package rdfiotypes

import (
	"context"
	"errors"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

var ErrResourceNotSupported = errors.New("resource not supported")

//

type ResourceMetadata interface {
	GetIRI() rdf.IRI
	GetFileName() (string, bool)
}

//

type ReaderOptions struct {
	Name   string
	Params []string
	Tee    io.Writer
}

//

type Reader interface {
	ResourceMetadata
	io.ReadCloser

	GetMediaType() (encoding.ContentMediaType, bool)
	GetMagicBytes() ([]byte, bool)

	AddTee(w io.Writer)
}

//

type WriterOptions struct {
	Name   string
	Params []string
	Tee    io.Writer
}

//

type Writer interface {
	ResourceMetadata
	io.WriteCloser
}

//

type ResourceManager interface {
	OpenReader(ctx context.Context, ref ReaderOptions) (Reader, error)
	OpenWriter(ctx context.Context, ref WriterOptions) (Writer, error)
}

//

type ResourceManagerList []ResourceManager

var _ ResourceManager = ResourceManagerList{}

func (rl ResourceManagerList) OpenReader(ctx context.Context, opts ReaderOptions) (Reader, error) {
	for _, rm := range rl {
		rr, err := rm.OpenReader(ctx, opts)
		if err == ErrResourceNotSupported {
			continue
		} else if err != nil {
			return nil, err
		}

		return rr, nil
	}

	return nil, ErrResourceNotSupported
}

func (rl ResourceManagerList) OpenWriter(ctx context.Context, opts WriterOptions) (Writer, error) {
	for _, rm := range rl {
		ww, err := rm.OpenWriter(ctx, opts)
		if err == ErrResourceNotSupported {
			continue
		} else if err != nil {
			return nil, err
		}

		return ww, nil
	}

	return nil, ErrResourceNotSupported
}
