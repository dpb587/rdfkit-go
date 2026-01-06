package encodingref

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

type ResourceWriter interface {
	io.WriteCloser
	ResourceMetadata
}

//

type ResourceReader interface {
	io.ReadCloser
	ResourceMetadata

	GetMediaType() (encoding.ContentMediaType, bool)
	GetMagicBytes() ([]byte, bool)

	AddTee(w io.Writer)
}

//

type ResourceRef struct {
	Name  string
	Flags []string
}

//

type ResourceManager interface {
	OpenWriter(ctx context.Context, ref ResourceRef) (ResourceWriter, error)
	OpenReader(ctx context.Context, ref ResourceRef) (ResourceReader, error)
}

//

type resourceManagerList []ResourceManager

func NewResourceManager(rms ...ResourceManager) ResourceManager {
	return resourceManagerList(rms)
}

func (rml resourceManagerList) OpenReader(ctx context.Context, ref ResourceRef) (ResourceReader, error) {
	for _, rm := range rml {
		rr, err := rm.OpenReader(ctx, ref)
		if err == ErrResourceNotSupported {
			continue
		} else if err != nil {
			return nil, err
		}

		return rr, nil
	}

	return nil, ErrResourceNotSupported
}

func (rml resourceManagerList) OpenWriter(ctx context.Context, ref ResourceRef) (ResourceWriter, error) {
	for _, rm := range rml {
		ww, err := rm.OpenWriter(ctx, ref)
		if err == ErrResourceNotSupported {
			continue
		} else if err != nil {
			return nil, err
		}

		return ww, nil
	}

	return nil, ErrResourceNotSupported
}
