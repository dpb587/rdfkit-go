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
	NewReader(ctx context.Context, ref ReaderOptions) (Reader, error)
	NewWriter(ctx context.Context, ref WriterOptions) (Writer, error)
}
