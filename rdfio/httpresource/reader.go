package httpresource

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type reader struct {
	iri          rdf.IRI
	body         io.Reader
	bodyCloser   func() error
	bodyPeek     []byte
	bodyPeekOnce sync.Once
	headers      http.Header
	encoding     string
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

	return rr.bodyCloser()
}

func (rr *reader) GetIRI() rdf.IRI {
	return rr.iri
}

func (rr *reader) GetFileName() (string, bool) {
	var filename string

	if headerValue := rr.headers.Get("Content-Disposition"); len(headerValue) > 0 {
		_, params, err := mime.ParseMediaType(headerValue)
		if err != nil {
			return "", false
		}

		filename = params["filename"]
	}

	if len(filename) == 0 {
		filename = filepath.Base(string(rr.iri))
	}

	if len(filename) == 0 {
		return "", false
	}

	return filename, true
}

func (rr *reader) GetMediaType() (encoding.ContentMediaType, bool) {
	headerValue := rr.headers.Get("Content-Type")
	if len(headerValue) == 0 {
		return encoding.ContentMediaType{}, false
	}

	mediatype, params, err := mime.ParseMediaType(headerValue)
	if err != nil {
		return encoding.ContentMediaType{}, false
	}

	parts := strings.SplitN(mediatype, "/", 2)
	if len(parts) != 2 {
		return encoding.ContentMediaType{}, false
	}

	return encoding.ContentMediaType{
		Type:       parts[0],
		Subtype:    parts[1],
		Parameters: params,
	}, true
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

func (rr *reader) GetOriginHeaders() http.Header {
	return rr.headers
}
