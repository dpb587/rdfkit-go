package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/htmljsonld"
	"github.com/dpb587/rdfkit-go/encoding/htmlmicrodata"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/rdfa"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
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

	htmlDocument, err := html.ParseDocument(
		res.Body,
		html.DocumentConfig{}.
			SetLocation(res.Request.URL.String()),
	)
	if err != nil {
		panic(fmt.Errorf("parse html: %v", err))
	}

	// Create all the HTML-based decoders since the original encoding is unknown.
	// (this might be simplified into a single, wrapper decoder in the future)

	htmlJsonld, err := htmljsonld.NewDecoder(
		htmlDocument,
		htmljsonld.DecoderConfig{}.
			SetDocumentLoader(jsonldtype.NewCachingDocumentLoader(
				jsonldtype.NewDefaultDocumentLoader(http.DefaultClient),
			)),
	)
	if err != nil {
		panic(fmt.Errorf("prepare htmljsonld: %v", err))
	}

	htmlMicrodata, err := htmlmicrodata.NewDecoder(
		htmlDocument,
		htmlmicrodata.DecoderConfig{}.
			SetVocabularyResolver(htmlmicrodata.ItemtypeVocabularyResolver),
	)
	if err != nil {
		panic(fmt.Errorf("prepare htmlmicrodata: %v", err))
	}

	htmlRdfa, err := rdfa.NewDecoder(htmlDocument)
	if err != nil {
		panic(fmt.Errorf("prepare rdfa: %v", err))
	}

	decoder := rdfioutil.NewStatementIteratorIterator(htmlJsonld, htmlMicrodata, htmlRdfa)

	defer decoder.Close()

	// Collect all the statements into subject-grouped statements.

	resourcesBuilder := rdfdescription.NewResourceListBuilder()

	for decoder.Next() {
		resourcesBuilder.AddTriple(decoder.GetStatement().GetTriple())
	}

	if err := decoder.Err(); err != nil {
		panic(fmt.Errorf("decode: %v", err))
	}

	// Prepare the Turtle encoder.

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

	// Output everything as resources since Turtle supports nested resource syntax.

	for _, resource := range resourcesBuilder.GetResources() {
		if err := encoder.PutResource(ctx, resource); err != nil {
			panic(fmt.Errorf("encode: %v", err))
		}
	}
}
