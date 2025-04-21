package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

func main() {
	parsed, err := inspectjson.Parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	expanded, err := jsonldinternal.Expand(
		parsed,
		jsonldtype.ProcessorOptions{
			BaseURL:        "https://stdin.local/",
			DocumentLoader: jsonldtype.NewDefaultDocumentLoader(http.DefaultClient),
		},
	)
	if err != nil {
		panic(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")

	err = e.Encode(expanded.AsBuiltin())
	if err != nil {
		panic(err)
	}
}
