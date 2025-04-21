package sparql

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory"
)

type Client struct {
	upstream *http.Client
	baseURL  string
}

func NewClient(upstream *http.Client, baseURL string) *Client {
	return &Client{
		upstream: upstream,
		baseURL:  baseURL,
	}
}

func (c *Client) Query(ctx context.Context, query string) (*QueryResponse, error) {
	formValues := url.Values{}
	formValues.Set("query", query)
	formEncoded := formValues.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL,
		bytes.NewBufferString(formEncoded),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEncoded)))
	req.Header.Set("Accept", "application/sparql-results+json")

	res, err := c.upstream.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return DecodeQueryResponse(res)
}

func (c *Client) Construct(ctx context.Context, query string) (rdfio.GraphStatementIterator, error) {
	formValues := url.Values{}
	formValues.Set("query", query)
	formEncoded := formValues.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL,
		bytes.NewBufferString(formEncoded),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEncoded)))
	req.Header.Set("Accept", "text/turtle")

	res, err := c.upstream.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	repo := inmemory.NewDataset()

	reader, err := turtle.NewDecoder(res.Body)
	if err != nil {
		return nil, fmt.Errorf("decoder: %v", err)
	}

	for reader.Next() {
		err := repo.PutTriple(ctx, reader.GetTriple())
		if err != nil {
			return nil, err
		}
	}

	if err := reader.Err(); err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return repo.NewStatementIterator(ctx) // TODO leaking repo.Close
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.upstream.Do(req)
}
