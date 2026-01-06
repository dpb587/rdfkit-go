package encodingref

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
)

type fileResourceManager struct{}

var _ ResourceManager = &fileResourceManager{}

func NewFileResourceManager() ResourceManager {
	return &fileResourceManager{}
}

func (rm *fileResourceManager) OpenReader(ctx context.Context, ref ResourceRef) (ResourceReader, error) {
	fp := strings.TrimPrefix(ref.Name, "file://")

	if fp == "-" {
		return &fileResourceReader{
			iri:      rdf.IRI("file:///dev/stdin"),
			fileName: "stdin",
			body:     os.Stdin,
		}, nil
	}

	f, err := os.OpenFile(fp, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return &fileResourceReader{
		iri:        rdf.IRI("file://" + fp),
		fileName:   filepath.Base(fp),
		body:       f,
		bodyCloser: f.Close,
	}, nil
}

func (rm *fileResourceManager) OpenWriter(ctx context.Context, ref ResourceRef) (ResourceWriter, error) {
	fp := strings.TrimPrefix(ref.Name, "file://")

	if fp == "-" {
		return &fileResourceWriter{
			iri:      rdf.IRI("file:///dev/stdout"),
			fileName: "stdout",
			body:     os.Stdout,
		}, nil
	}

	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return &fileResourceWriter{
		iri:        rdf.IRI("file://" + fp),
		fileName:   filepath.Base(fp),
		body:       f,
		bodyCloser: f.Close,
	}, nil
}

//

type fileResourceReader struct {
	iri          rdf.IRI
	fileName     string
	body         io.Reader
	bodyCloser   func() error
	bodyPeek     []byte
	bodyPeekOnce sync.Once
}

var _ ResourceReader = &fileResourceReader{}

func (rr *fileResourceReader) AddTee(w io.Writer) {
	rr.body = io.TeeReader(rr.body, w)
}

func (rr *fileResourceReader) Read(p []byte) (n int, err error) {
	return rr.body.Read(p)
}

func (rr *fileResourceReader) Close() error {
	rr.body = nil // panic if further reads are attempted

	if rr.bodyCloser == nil {
		return nil
	}

	return rr.bodyCloser()
}

func (rr *fileResourceReader) GetIRI() rdf.IRI {
	return rr.iri
}

func (rr *fileResourceReader) GetFileName() (string, bool) {
	return rr.fileName, len(rr.fileName) > 0
}

func (rr *fileResourceReader) GetMediaType() (encoding.ContentMediaType, bool) {
	return encoding.ContentMediaType{}, false
}

func (rr *fileResourceReader) GetMagicBytes() ([]byte, bool) {
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

//

type fileResourceWriter struct {
	iri        rdf.IRI
	fileName   string
	body       io.Writer
	bodyCloser func() error
}

var _ ResourceWriter = &fileResourceWriter{}

func (ww *fileResourceWriter) Write(p []byte) (n int, err error) {
	return ww.body.Write(p)
}

func (ww *fileResourceWriter) Close() error {
	ww.body = nil // panic if further writes are attempted

	if ww.bodyCloser == nil {
		return nil
	}

	return ww.bodyCloser()
}

func (ww *fileResourceWriter) GetIRI() rdf.IRI {
	return ww.iri
}

func (ww *fileResourceWriter) GetFileName() (string, bool) {
	return ww.fileName, len(ww.fileName) > 0
}
