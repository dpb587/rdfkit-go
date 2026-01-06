package fileresource

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type manager struct{}

var _ rdfiotypes.ResourceManager = &manager{}

func NewManager() rdfiotypes.ResourceManager {
	return &manager{}
}

func (e manager) NewReaderParams() rdfiotypes.Params {
	return &readerParams{}
}

func (rm *manager) NewReader(ctx context.Context, opts rdfiotypes.ReaderOptions) (rdfiotypes.Reader, error) {
	flags := &readerParams{}

	err := rdfiotypes.LoadAndApplyParams(flags, opts.Params...)
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

	return rr, nil
}

func (rm *manager) NewWriter(ctx context.Context, opts rdfiotypes.WriterOptions) (rdfiotypes.Writer, error) {
	fp := strings.TrimPrefix(opts.Name, "file://")

	var ww *writer

	if fp == "" || fp == "-" {
		ww = &writer{
			iri:      rdf.IRI("file:///dev/stdout"),
			fileName: "stdout",
			body:     os.Stdout,
		}
	} else {
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}

		ww = &writer{
			iri:        rdf.IRI("file://" + fp),
			fileName:   filepath.Base(fp),
			body:       f,
			bodyCloser: f.Close,
		}
	}

	if opts.Tee != nil {
		ww.body = io.MultiWriter(ww.body, opts.Tee)
	}

	return ww, nil
}

//

//
