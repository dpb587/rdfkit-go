package jsonldtype

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/tomnomnom/linkheader"
)

type DefaultDocumentLoader struct {
	client *http.Client
}

var _ DocumentLoader = &DefaultDocumentLoader{}

func NewDefaultDocumentLoader(client *http.Client) *DefaultDocumentLoader {
	dl := &DefaultDocumentLoader{
		client: client,
	}

	if dl.client == nil {
		dl.client = http.DefaultClient
	}

	return dl
}

func (dl *DefaultDocumentLoader) LoadDocument(ctx context.Context, u string, opts DocumentLoaderOptions) (RemoteDocument, error) {
	var maxRequests = 10
	var nextURL *url.URL
	var err error

	nextURL, err = url.Parse(u)
	if err != nil {
		return RemoteDocument{}, err
	}

	for {
		if maxRequests == 0 {
			return RemoteDocument{}, fmt.Errorf("fetch: maximum number of requests reached")
		}

		maxRequests--

		// normalize

		if nextURL.Fragment != "" {
			nextURL.Fragment = ""
		}

		//

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL.String(), nil)
		if err != nil {
			return RemoteDocument{}, err
		}

		{
			var acceptValues = []string{
				"application/ld+json",
				"application/json;q=0.9",
				"*/*;q=0.1",
			}

			switch len(opts.RequestProfile) {
			case 0:
				// nop
			case 1:
				acceptValues[0] += ";profile=" + opts.RequestProfile[0]
			default:
				acceptValues[0] += `;profile="` + strings.Join(opts.RequestProfile, " ") + `"`
			}

			req.Header.Set("Accept", strings.Join(acceptValues, ", "))
		}

		req.Header.Set("Accept-Encoding", "gzip, deflate")

		resp, err := dl.client.Do(req)
		if err != nil {
			return RemoteDocument{}, fmt.Errorf("fetch: %w", err)
		}

		if resp.Body != nil {
			defer resp.Body.Close()
		}

		// [spec // 3] Set *documentUrl* to the location of the retrieved resource considering redirections (exclusive of HTTP status `303` "See Other" redirects as discussed in [cooluris]).

		documentUrl := resp.Request.URL

		// [spec // 4] If the retrieved resource's Content-Type is not `application/json` nor any media type with a `+json` suffix as defined in [RFC6839], and the response has an HTTP Link Header [RFC8288] using the `alternate` link relation with type `application/ld+json`, set *url* to the associated `href` relative to the previous *url* and restart the algorithm from step 2.

		respContentType := strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)[0]
		respLinks := linkheader.ParseMultiple(resp.Header["Link"])

		if respContentType != "application/ld+json" && !strings.HasSuffix(respContentType, "+json") {
			var followLink *linkheader.Link

			for _, link := range respLinks {
				if link.Rel == "alternate" && link.Params["type"] == "application/ld+json" {
					followLink = &link

					break
				}
			}

			if followLink != nil {
				nextURL, err = documentUrl.Parse(followLink.URL)
				if err != nil {
					return RemoteDocument{}, fmt.Errorf("failed to parse alternate link: %v", err)
				}

				continue
			}
		}

		// [dpb] status code not explicitly listed in spec; currently following links on all status codes, but otherwise failing

		if resp.StatusCode != http.StatusOK {
			return RemoteDocument{}, Error{
				Code: LoadingDocumentFailed,
				Err:  fmt.Errorf("unexpected status code: %d", resp.StatusCode),
			}
		}

		// [spec // 5] If the retrieved resource's Content-Type is `application/json` or any media type with a `+json` suffix as defined in [RFC6839] except `application/ld+json`, and the response has an HTTP Link Header [RFC8288] using the `http://www.w3.org/ns/json-ld#context` link relation, set *contextUrl* to the associated `href`.
		// [spec // 5] If multiple HTTP Link Headers using the `http://www.w3.org/ns/json-ld#context` link relation are found, the promise is rejected with a `JsonLdError` whose `code` is set to `multiple context link headers` and processing is terminated.
		// [spec // 5] Processors *MAY* transform *document* to the internal representation.
		// [spec // 5] NOTE The HTTP Link Header is ignored for documents served as `application/ld+json`, `text/html`, or `application/xhtml+xml`.

		switch respContentType {
		case "application/ld+json", "text/html", "application/xhtml+xml":
			respLinks = nil
		}

		var contextUrl *url.URL

		if respContentType == "application/json" || strings.HasSuffix(respContentType, "+json") {
			for _, link := range respLinks {
				if link.Rel == "http://www.w3.org/ns/json-ld#context" {
					if contextUrl != nil {
						return RemoteDocument{}, Error{
							Code: MultipleContextLinkHeaders,
						}
					}

					contextUrl, err = documentUrl.Parse(link.URL)
					if err != nil {
						return RemoteDocument{}, fmt.Errorf("failed to parse content link: %v", err)
					}
				}
			}
		} else {

			// [spec // 6] Otherwise, the retrieved document's Content-Type is neither `application/json`, `application/ld+json`, nor any other media type using a `+json` suffix as defined in [RFC6839]. Reject the promise passing a `loading document failed` error.

			return RemoteDocument{}, Error{
				Code: LoadingDocumentFailed,
				Err:  fmt.Errorf("unexpected content type: %s", respContentType),
			}
		}

		// [spec // 7] Create a new `RemoteDocument` *remote document* using *url* as `documentUrl`, *document* as `document`, the returned Content-Type (without parameters) as `contentType`, any returned `profile` parameter, or `null` as `profile`, and *contextUrl*, or `null` as *contextUrl*.

		var bodyReader io.Reader = resp.Body

		switch contentEncoding := resp.Header.Get("Content-Encoding"); contentEncoding {
		case "":
			// nop
		case "gzip":
			bodyReader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return RemoteDocument{}, Error{
					Code: LoadingDocumentFailed,
					Err:  fmt.Errorf("decode gzip: %v", err),
				}
			}
		case "deflate":
			bodyReader = flate.NewReader(resp.Body)
		default:
			return RemoteDocument{}, Error{
				Code: LoadingDocumentFailed,
				Err:  fmt.Errorf("unsupported content encoding: %s", contentEncoding),
			}
		}

		parsed, err := inspectjson.Parse(bodyReader)
		if err != nil {
			return RemoteDocument{}, Error{
				Code: LoadingDocumentFailed,
				Err:  fmt.Errorf("parse (%s): %v", documentUrl, err),
			}
		}

		return RemoteDocument{
			DocumentURL: documentUrl,
			Document:    parsed,
			ContentType: respContentType,
			Profile:     resp.Header.Get("Profile"),
			ContextURL:  contextUrl,
		}, nil
	}
}
