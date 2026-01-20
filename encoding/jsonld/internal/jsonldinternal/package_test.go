package jsonldinternal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

func TestOne(t *testing.T) {

	parsed, err := inspectjson.Parse(
		strings.NewReader(`{
  "@context": {
    "@base": "http://example.org/"
  },
  "http://example.org/vocab/at": {"@id": "@"},
  "http://example.org/vocab/foo.bar": {"@id": "@foo.bar"},
  "http://example.org/vocab/ignoreme": {"@id": "@ignoreMe"}
}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expanded, err := Expand(
		parsed,
		jsonldtype.ProcessorOptions{
			BaseURL:        "http://units.example.com/sub/path",
			DocumentLoader: jsonldtype.NewDefaultDocumentLoader(http.DefaultClient),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := json.NewEncoder(os.Stderr)
	e.SetIndent("", "  ")
	e.Encode(expanded.AsBuiltin())
}

func TestBugIgnoredBaseAfterRemoteContext(t *testing.T) {
	parsed, err := inspectjson.Parse(strings.NewReader(`
		{
			"@context": [
				"context.jsonld",
				{
					"@base": "toRdf-manifest"
				}
			],
			"@id": "#t0001",
			"name": "test"
		}
	`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expanded, err := Expand(
		parsed,
		jsonldtype.ProcessorOptions{
			BaseURL: "http://example.com/tests/toRdf-manifest.jsonld",
			DocumentLoader: jsonldtype.DocumentLoaderFunc(func(ctx context.Context, urlv string, opts jsonldtype.DocumentLoaderOptions) (jsonldtype.RemoteDocument, error) {
				u, _ := url.Parse(urlv)

				return jsonldtype.RemoteDocument{
					ContentType: "application/ld+json",
					Document: inspectjson.ObjectValue{
						Members: map[string]inspectjson.ObjectMember{
							"@context": {
								Value: inspectjson.ObjectValue{
									Members: map[string]inspectjson.ObjectMember{
										"@vocab": {
											Value: inspectjson.StringValue{Value: "http://example.org/"},
										},
									},
								},
							},
						},
					},
					DocumentURL: u,
				}, nil
			}),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := json.Marshal(expanded.AsBuiltin())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e := "http://example.com/tests/toRdf-manifest#t0001"; !strings.Contains(string(result), _e) {
		t.Fatalf("expected %q, but found none", _e)
	}
}
