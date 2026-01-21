package fileresource

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/rdfkit-go/internal/ioutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

var WriterBufferSize = 32 * 1024

type manager struct{}

var _ rdfiotypes.ResourceManager = &manager{}

func NewManager() rdfiotypes.ResourceManager {
	return &manager{}
}

func (e manager) NewReaderParams() rdfiotypes.Params {
	return newReaderParams()
}

func (rm *manager) NewReader(ctx context.Context, opts rdfiotypes.ReaderOptions) (rdfiotypes.Reader, error) {
	params := newReaderParams()

	err := rdfiotypes.LoadAndApplyParams(params, opts.Params...)
	if err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	fp := strings.TrimPrefix(opts.Name, "file://")

	var rr *reader

	if fp == "" || fp == "-" {
		rr = &reader{
			iri:      rdf.IRI("file:///dev/stdin"),
			fileName: "stdin",
			body:     os.Stdin,
		}
	} else {
		f, err := os.OpenFile(fp, os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}

		rr = &reader{
			iri:        rdf.IRI("file://" + fp),
			fileName:   filepath.Base(fp),
			body:       f,
			bodyCloser: f.Close,
		}
	}

	if opts.Tee != nil {
		rr.body = io.TeeReader(rr.body, opts.Tee)
	}

	err = params.Filter.ResolveReader(
		rr.body,
		rr.bodyCloser,
		func() []string {
			if filename, ok := rr.GetFileName(); ok {
				switch strings.ToLower(filepath.Ext(filename)) {
				case ".gz":
					return []string{"gzip"}
				}
			}

			// fallback to all

			return []string{"gzip"}
		},
		func(nextReader io.Reader, nextCloser ioutil.CloserFunc) {
			rr.body = nextReader
			rr.bodyCloser = nextCloser
		},
	)
	if err != nil {
		rr.bodyCloser()

		return nil, err
	}

	return rr, nil
}

func (e manager) NewWriterParams() rdfiotypes.Params {
	return &writerParams{}
}

func (rm *manager) NewWriter(ctx context.Context, opts rdfiotypes.WriterOptions) (rdfiotypes.Writer, error) {
	params := newWriterParams()

	err := rdfiotypes.LoadAndApplyParams(params, opts.Params...)
	if err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	fp := strings.TrimPrefix(opts.Name, "file://")

	var ww *writer

	if fp == "" || fp == "-" {
		bufWriter := bufio.NewWriterSize(os.Stdout, WriterBufferSize)

		ww = &writer{
			iri:      rdf.IRI("file:///dev/stdout"),
			fileName: "stdout",
			body:     bufWriter,
			bodyCloser: func() error {
				return bufWriter.Flush()
			},
		}
	} else {
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}

		bufWriter := bufio.NewWriterSize(f, WriterBufferSize)

		ww = &writer{
			iri:      rdf.IRI("file://" + fp),
			fileName: filepath.Base(fp),
			body:     bufWriter,
			bodyCloser: func() error {
				err := bufWriter.Flush()
				if err != nil {
					return err
				}

				return f.Close()
			},
		}
	}

	if opts.Tee != nil {
		ww.body = io.MultiWriter(ww.body, opts.Tee)
	}

	err = params.Filter.ResolveWriterEncoding(
		ww.body,
		ww.bodyCloser,
		func() (string, bool) {
			if filename, ok := ww.GetFileName(); ok {
				switch strings.ToLower(filepath.Ext(filename)) {
				case ".gz":
					return "gzip", true
				}
			}

			return "", false
		},
		func(nextWriter io.Writer, nextCloser ioutil.CloserFunc) {
			ww.body = nextWriter
			ww.bodyCloser = nextCloser
		},
	)
	if err != nil {
		ww.bodyCloser()

		return nil, err
	}

	return ww, nil
}
