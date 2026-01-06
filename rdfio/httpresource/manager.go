package httpresource

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dpb587/rdfkit-go/internal/ioutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type manager struct {
	client *http.Client
}

var _ rdfiotypes.ResourceManager = &manager{}

func NewManager(client *http.Client) rdfiotypes.ResourceManager {
	return &manager{
		client: client,
	}
}

func (e manager) NewReaderParams() rdfiotypes.Params {
	return newReaderParams()
}

func (rm *manager) NewReader(ctx context.Context, opts rdfiotypes.ReaderOptions) (rdfiotypes.Reader, error) {
	if !strings.HasPrefix(opts.Name, "http://") && !strings.HasPrefix(opts.Name, "https://") {
		return nil, rdfiotypes.ErrResourceNotSupported
	}

	params := newReaderParams()

	err := rdfiotypes.LoadAndApplyParams(params, opts.Params...)
	if err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", opts.Name, nil)
	if err != nil {
		return nil, rdfiotypes.ErrResourceNotSupported
	}

	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := rm.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.Body != nil {
			resp.Body.Close()
		}

		return nil, fmt.Errorf("web: unexpected status code: %d", resp.StatusCode)
	}

	rr := &reader{
		iri:        rdf.IRI(resp.Request.URL.String()),
		body:       resp.Body,
		bodyCloser: resp.Body.Close,
		headers:    resp.Header,
	}

	if encoding := strings.ToLower(resp.Header.Get("Content-Encoding")); len(encoding) > 0 {
		switch encoding {
		case "gzip":
			gunzip, err := gzip.NewReader(rr.body)
			if err != nil {
				rr.Close()

				return nil, fmt.Errorf("web: encoding[=gzip]: %w", err)
			}

			rr.body = gunzip
			rr.bodyCloser = func() error {
				err := gunzip.Close()
				if err != nil {
					return err
				}

				return resp.Body.Close()
			}
		default:
			return nil, fmt.Errorf("web: unexpected encoding %q", encoding)
		}
	}

	if opts.Tee != nil {
		rr.body = io.TeeReader(rr.body, opts.Tee)
	}

	err = params.Filter.ResolveReader(
		rr.body,
		rr.bodyCloser,
		func() []string {
			switch strings.ToLower(strings.SplitN(rr.headers.Get("Content-Type"), ";", 2)[0]) {
			case "":
				// unknown
			case "application/octet-stream":
				// possible
			default:
				return nil
			}

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
		resp.Body.Close()

		return nil, err
	}

	return rr, nil
}

func (rm *manager) NewWriter(ctx context.Context, opts rdfiotypes.WriterOptions) (rdfiotypes.Writer, error) {
	return nil, rdfiotypes.ErrResourceNotSupported
}
