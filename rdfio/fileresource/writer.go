package fileresource

import (
	"io"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type writer struct {
	iri        rdf.IRI
	fileName   string
	body       io.Writer
	bodyCloser func() error
}

var _ rdfiotypes.Writer = &writer{}

func (ww *writer) Write(p []byte) (n int, err error) {
	return ww.body.Write(p)
}

func (ww *writer) Close() error {
	ww.body = nil // panic if further writes are attempted

	if ww.bodyCloser == nil {
		return nil
	}

	return ww.bodyCloser()
}

func (ww *writer) GetIRI() rdf.IRI {
	return ww.iri
}

func (ww *writer) GetFileName() (string, bool) {
	return ww.fileName, len(ww.fileName) > 0
}
