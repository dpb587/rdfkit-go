package fileresource

import (
	"bytes"
	"io"
	"sync"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type reader struct {
	iri          rdf.IRI
	fileName     string
	body         io.Reader
	bodyCloser   func() error
	bodyPeek     []byte
	bodyPeekOnce sync.Once
}

var _ rdfiotypes.Reader = &reader{}

func (rr *reader) AddTee(w io.Writer) {
	rr.body = io.TeeReader(rr.body, w)
}

func (rr *reader) Read(p []byte) (n int, err error) {
	return rr.body.Read(p)
}

func (rr *reader) Close() error {
	rr.body = nil // panic if further reads are attempted

	if rr.bodyCloser == nil {
		return nil
	}

	return rr.bodyCloser()
}

func (rr *reader) GetIRI() rdf.IRI {
	return rr.iri
}

func (rr *reader) GetFileName() (string, bool) {
	return rr.fileName, len(rr.fileName) > 0
}

func (rr *reader) GetMediaType() (encoding.ContentMediaType, bool) {
	return encoding.ContentMediaType{}, false
}

func (rr *reader) GetMagicBytes() ([]byte, bool) {
	rr.bodyPeekOnce.Do(func() {
		var buf = make([]byte, 1024)

		l, _ := io.ReadFull(rr.body, buf)
		// assumes err will be propagated by later reads?

		rr.bodyPeek = buf[:l]

		if l > 0 {
			rr.body = io.MultiReader(bytes.NewReader(buf[:l]), rr.body)
		}
	})

	return rr.bodyPeek, len(rr.bodyPeek) > 0
}
