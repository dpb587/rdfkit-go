package jsonldtype

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dpb587/inspectjson-go/inspectjson"
)

type DocumentLoaderOptions struct {
	ExtractAllScripts bool
	Profile           *string
	RequestProfile    []string
}

type DocumentLoader interface {
	LoadDocument(ctx context.Context, url string, opts DocumentLoaderOptions) (RemoteDocument, error)
}

type DocumentLoaderFunc func(ctx context.Context, url string, opts DocumentLoaderOptions) (RemoteDocument, error)

func (f DocumentLoaderFunc) LoadDocument(ctx context.Context, url string, opts DocumentLoaderOptions) (RemoteDocument, error) {
	return f(ctx, url, opts)
}

// [spec // 9.4.3] The RemoteDocument type is used by a LoadDocumentCallback to return information about a remote document or context.
type RemoteDocument struct {
	ContentType string
	ContextURL  *url.URL
	Document    inspectjson.Value
	DocumentURL *url.URL
	Profile     string
	Headers     http.Header
}
