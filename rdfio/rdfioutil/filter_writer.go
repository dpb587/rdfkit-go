package rdfioutil

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/dpb587/kvstrings-go/kvstrings"
	"github.com/dpb587/kvstrings-go/kvstrings/kvref"
	"github.com/dpb587/rdfkit-go/internal/ioutil"
	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type FilterWriter struct {
	Filter          *string
	FilterGzipLevel *int
}

func (f *FilterWriter) NewParamsCollection(base kvstrings.KeyName) rdfiotypes.ParamsCollection {
	return rdfiotypes.ParamsCollection{
		base: kvref.StringPtr(&f.Filter, rdfiotypes.ParamMeta{
			Usage: "Filter output (one of auto, none, gzip; default auto)",
		}),
		base + ".gzip.level": kvref.IntPtr(&f.FilterGzipLevel, rdfiotypes.ParamMeta{
			Usage: "Set the gzip compression level (requires filter=gzip)",
		}),
	}
}

func (f *FilterWriter) ApplyDefaults() {
	if f.Filter == nil {
		f.Filter = ptr.Value("auto")
	}

	if *f.Filter == "auto" || *f.Filter == "gzip" {
		f.FilterGzipLevel = ptr.Value(gzip.DefaultCompression)
	}
}

func (f *FilterWriter) ResolveWriterEncoding(
	originalWriter io.Writer,
	originalCloser ioutil.CloserFunc,
	getAutoEncoding func() (string, bool),
	updateFunc func(nextWriter io.Writer, nextCloser ioutil.CloserFunc),
) error {
	if *f.Filter == "none" {
		return nil
	}

	encoding := *f.Filter

	if encoding == "auto" {
		var ok bool

		encoding, ok = getAutoEncoding()
		if !ok {
			return nil
		}
	}

	switch encoding {
	case "gzip":
		gzipWriter, err := gzip.NewWriterLevel(originalWriter, *f.FilterGzipLevel)
		if err != nil {
			return fmt.Errorf("filter[=gzip]: %w", err)
		}

		updateFunc(gzipWriter, func() error {
			err := gzipWriter.Close()
			if err != nil {
				return err
			}

			return originalCloser()
		})

		return nil
	}

	return fmt.Errorf("unknown filter %q", *f.Filter)
}
