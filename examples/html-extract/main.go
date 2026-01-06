package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dpb587/rdfkit-go/encoding/html/htmldefaults"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		panic("Usage: html-extract URL")
	}

	res, err := http.DefaultClient.Get(os.Args[1])
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected response status: %s", res.Status))
	}

	// Parse the HTML document into a DOM tree.

	decoder, err := htmldefaults.NewDecoder(
		res.Body,
		htmldefaults.DecoderConfig{}.
			SetLocation(res.Request.URL.String()).
			SetDocumentLoaderJSONLD(jsonldtype.NewCachingDocumentLoader(
				jsonldtype.NewDefaultDocumentLoader(http.DefaultClient),
			)),
	)
	if err != nil {
		panic(fmt.Errorf("creating decoder: %v", err))
	}

	defer decoder.Close()

	// Load all statements.

	resourceEncoder := rdfdescription.NewResourceListBuilder()

	for decoder.Next() {
		// ignoring the graph name for now
		resourceEncoder.AddTriple(ctx, decoder.Quad().Triple)
	}

	if err := decoder.Err(); err != nil {
		panic(fmt.Errorf("decode: %v", err))
	}

	// Output as structured Turtle resources.

	encoder, err := turtle.NewEncoder(
		os.Stdout,
		turtle.EncoderConfig{}.
			SetBase(res.Request.URL.String()).
			SetPrefixes(iriutil.NewPrefixMap(rdfacontext.WidelyUsedInitialContext()...)).
			SetBuffered(true),
	)
	if err != nil {
		panic(fmt.Errorf("turtle encoder: %v", err))
	}

	defer encoder.Close()

	err = resourceEncoder.AddTo(ctx, encoder, true)
	if err != nil {
		panic(fmt.Errorf("encode: %v", err))
	}
}
