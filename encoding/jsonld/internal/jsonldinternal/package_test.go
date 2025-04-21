package jsonldinternal

import (
	"encoding/json"
	"net/http"
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
