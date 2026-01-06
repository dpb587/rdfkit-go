package httpresource

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	return &readerParams{}
}

func (rm *manager) NewReader(ctx context.Context, opts rdfiotypes.ReaderOptions) (rdfiotypes.Reader, error) {
	if !strings.HasPrefix(opts.Name, "http://") && !strings.HasPrefix(opts.Name, "https://") {
		return nil, rdfiotypes.ErrResourceNotSupported
	}

	flags := &readerParams{}

	err := rdfiotypes.LoadAndApplyParams(flags, opts.Params...)
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

	return rr, nil
}

func (rm *manager) NewWriter(ctx context.Context, opts rdfiotypes.WriterOptions) (rdfiotypes.Writer, error) {
	return nil, rdfiotypes.ErrResourceNotSupported
}
