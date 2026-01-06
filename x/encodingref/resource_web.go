package encodingref

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type webResourceManager struct {
	client *http.Client
}

var _ ResourceManager = &webResourceManager{}

func NewWebResourceManager(client *http.Client) ResourceManager {
	return &webResourceManager{
		client: client,
	}
}

func (rm *webResourceManager) OpenReader(ctx context.Context, ref ResourceRef) (ResourceReader, error) {
	if !strings.HasPrefix(ref.Name, "http://") && !strings.HasPrefix(ref.Name, "https://") {
		return nil, ErrResourceNotSupported
	}

	req, err := http.NewRequestWithContext(ctx, "GET", ref.Name, nil)
	if err != nil {
		return nil, ErrResourceNotSupported
	}

	resp, err := rm.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.Body != nil {
			resp.Body.Close()
		}

		return nil, fmt.Errorf("web: unexpected status code: %d", resp.StatusCode)
	}

	return &webResourceReader{
		iri:        rdf.IRI(resp.Request.URL.String()),
		body:       resp.Body,
		bodyCloser: resp.Body.Close,
		headers:    resp.Header,
	}, nil
}

func (rm *webResourceManager) OpenWriter(ctx context.Context, ref ResourceRef) (ResourceWriter, error) {
	return nil, ErrResourceNotSupported
}

//

type webResourceReader struct {
	iri          rdf.IRI
	body         io.Reader
	bodyCloser   func() error
	bodyPeek     []byte
	bodyPeekOnce sync.Once
	headers      http.Header
}

var _ ResourceReader = &webResourceReader{}

func (rr *webResourceReader) AddTee(w io.Writer) {
	rr.body = io.TeeReader(rr.body, w)
}

func (rr *webResourceReader) Read(p []byte) (n int, err error) {
	return rr.body.Read(p)
}

func (rr *webResourceReader) Close() error {
	rr.body = nil // panic if further reads are attempted

	return rr.bodyCloser()
}

func (rr *webResourceReader) GetIRI() rdf.IRI {
	return rr.iri
}

func (rr *webResourceReader) GetFileName() (string, bool) {
	headerValue := rr.headers.Get("Content-Disposition")
	if len(headerValue) == 0 {
		return "", false
	}

	_, params, err := mime.ParseMediaType(headerValue)
	if err != nil {
		return "", false
	}

	filename, ok := params["filename"]
	return filename, ok
}

func (rr *webResourceReader) GetMediaType() (encoding.ContentMediaType, bool) {
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

func (rr *webResourceReader) GetMagicBytes() ([]byte, bool) {
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

func (rr *webResourceReader) GetOriginHeaders() http.Header {
	return rr.headers
}
