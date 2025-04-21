package jsonldtype

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type cachingDocumentLoaderResult struct {
	RemoteDocument RemoteDocument
	Err            error
}

type CachingDocumentLoader struct {
	upstream DocumentLoader
	cache    map[string]*cachingDocumentLoaderResult
}

func NewCachingDocumentLoader(upstream DocumentLoader) *CachingDocumentLoader {
	return &CachingDocumentLoader{
		upstream: upstream,
		cache:    map[string]*cachingDocumentLoaderResult{},
	}
}

var _ DocumentLoader = &CachingDocumentLoader{}

func (dl CachingDocumentLoader) LoadDocument(ctx context.Context, u string, opts DocumentLoaderOptions) (RemoteDocument, error) {
	cacheKey := dl.cacheKey(u, opts)

	if res, ok := dl.cache[cacheKey]; ok {
		return res.RemoteDocument, res.Err
	}

	remoteDocument, err := dl.upstream.LoadDocument(ctx, u, opts)
	cdlr := &cachingDocumentLoaderResult{
		RemoteDocument: remoteDocument,
		Err:            err,
	}

	dl.cache[cacheKey] = cdlr

	if cdlr.Err == nil && remoteDocument.DocumentURL.String() != u {
		dl.cache[dl.cacheKey(remoteDocument.DocumentURL.String(), opts)] = cdlr
	}

	return remoteDocument, err
}

func (dl CachingDocumentLoader) cacheKey(u string, opts DocumentLoaderOptions) string {
	cacheKeyHash := sha256.New()
	fmt.Fprintf(cacheKeyHash, "url: %s\nextractAllScripts: %v\n", u, opts.ExtractAllScripts)

	if opts.Profile != nil {
		fmt.Fprintf(cacheKeyHash, "profile: %s\n", *opts.Profile)
	}

	if len(opts.RequestProfile) > 0 {
		fmt.Fprintf(cacheKeyHash, "requestProfile: %s\n", strings.Join(opts.RequestProfile, " "))
	}

	return hex.EncodeToString(cacheKeyHash.Sum(nil))
}
