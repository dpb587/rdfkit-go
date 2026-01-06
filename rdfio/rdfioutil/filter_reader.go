package rdfioutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/dpb587/kvstrings-go/kvstrings"
	"github.com/dpb587/kvstrings-go/kvstrings/kvref"
	"github.com/dpb587/rdfkit-go/internal/ioutil"
	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type FilterReader struct {
	Filter *string
}

func (f *FilterReader) NewParamsCollection(base kvstrings.KeyName) rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{
		base: kvref.StringPtr(&f.Filter, rdfiotypes.ParamMeta{
			Usage: "Filter input (one of auto, none, gzip; default auto)",
		}),
	}
}

func (f *FilterReader) ApplyDefaults() {
	if f.Filter == nil {
		f.Filter = ptr.Value("auto")
	}
}

func (f *FilterReader) ResolveReader(
	originalReader io.Reader,
	originalCloser ioutil.CloserFunc,
	getAutoEncodings func() []string,
	updateFunc func(nextReader io.Reader, nextCloser ioutil.CloserFunc),
) error {
	if *f.Filter == "none" {
		return nil
	}

	var peek *bytes.Buffer

	doUnwrap := func(f func(r io.Reader) (io.Reader, ioutil.CloserFunc, error)) (io.Reader, ioutil.CloserFunc, error) {
		var unwrapReader io.Reader

		if peek != nil && peek.Len() > 0 {
			unwrapReader = io.MultiReader(bytes.NewReader(peek.Bytes()), originalReader)
		} else {
			unwrapReader = originalReader
		}

		peek = &bytes.Buffer{}

		return f(io.TeeReader(unwrapReader, peek))
	}

	unwrapGzip := func(r io.Reader) (io.Reader, ioutil.CloserFunc, error) {
		gzipReader, err := gzip.NewReader(r)
		if err != nil {
			return nil, nil, err
		}

		return gzipReader, func() error {
			// TODO multierr
			err := gzipReader.Close()
			originalCloser()

			return err
		}, nil
	}

	switch *f.Filter {
	case "auto":
		tries := getAutoEncodings()

		for _, unwrapTry := range tries {
			switch unwrapTry {
			case "gzip":
				if nextReader, nextCloser, err := doUnwrap(unwrapGzip); err == nil {
					updateFunc(nextReader, nextCloser)

					return nil
				}
			}
		}

		if peek != nil && peek.Len() > 0 {
			updateFunc(io.MultiReader(bytes.NewReader(peek.Bytes()), originalReader), originalCloser)

			return nil
		}

		return nil
	case "gzip":
		nextReader, nextCloser, err := doUnwrap(unwrapGzip)
		if err != nil {
			return fmt.Errorf("filter[=gzip]: %w", err)
		}

		updateFunc(nextReader, nextCloser)

		return nil
	}

	return fmt.Errorf("unknown filter %q", *f.Filter)
}
