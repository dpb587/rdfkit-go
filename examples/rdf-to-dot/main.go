package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdf/quads"
)

func main() {
	input := cmdflags.EncodingInput{
		FallbackOpener: cmdflags.WebRemoteOpener,
		DocumentLoaderJSONLD: jsonldtype.NewCachingDocumentLoader(
			jsonldtype.NewDefaultDocumentLoader(http.DefaultClient),
		),
	}

	flag.StringVar(&input.Path, "i", "-", "input path, url, or '-' for stdin")
	flag.StringVar(&input.Type, "t", "", "input format")
	flag.Parse()

	// Open the file/url and initialize the decoder.
	ih, err := input.Open()
	if err != nil {
		panic(err)
	}

	defer ih.Close()

	// Read all the statements before we visualize it. This is primarily to capture any base and prefix directives that
	// may be within the file.
	allStatements, err := quads.Collect(ih.Decoder)
	if err != nil {
		panic(err)
	}

	// Utilities for shortening how IRIs are displayed (based on directives from the input).

	var base *iriutil.BaseIRI

	if len(ih.DecodedBase) > 0 {
		// Last-most base ends up being used.
		base, err = iriutil.ParseBaseIRI(ih.DecodedBase[len(ih.DecodedBase)-1])
		if err != nil {
			panic(fmt.Errorf("decode base: %v", err))
		}
	}

	prefixes := iriutil.NewPrefixMap(append(
		// Include some well-known prefixes, like xsd.
		rdfacontext.WidelyUsedInitialContext(),
		// And override with prefixes from the input.
		ih.DecodedPrefixMappings...,
	)...)

	shortenIRI := func(iri rdf.IRI) string {
		if prefix, localName, ok := prefixes.CompactPrefix(iri); ok {
			return fmt.Sprintf("%s:%s", prefix, localName)
		}

		if base != nil {
			if relative, ok := base.RelativizeIRI(iri); ok {
				return fmt.Sprintf("<%s>", relative)
			}
		}

		return fmt.Sprintf("<%s>", iri)
	}

	// Utilities for adding nodes to the graph.

	literalNodes := 0
	knownResources := map[rdf.SubjectValue]string{}
	blankNodeStringer := blanknodeutil.NewStringerInt64()

	requireResource := func(w io.Writer, indent string, s rdf.SubjectValue) string {
		if node, ok := knownResources[s]; ok {
			return node
		}

		node := fmt.Sprintf("r%d", len(knownResources))
		knownResources[s] = node

		if sBlankNode, ok := s.(rdf.BlankNode); ok {
			fmt.Fprintf(w, indent+`%s [fillcolor="lavender",label=%q,shape=box,style="dashed,filled,rounded,setlinewidth(2)"]`+"\n", node, "_:"+blankNodeStringer.GetBlankNodeIdentifier(sBlankNode))
		} else if sIRI, ok := s.(rdf.IRI); ok {
			fmt.Fprintf(w, indent+`%s [fillcolor="lavender",href=%q,label=%q,shape=box,style="filled,rounded,setlinewidth(2)"]`+"\n", node, sIRI, shortenIRI(sIRI))
		} else {
			panic(fmt.Errorf("invalid subject type: %T", s))
		}

		return node
	}

	// Render the visualization.

	_, err = fmt.Fprintf(os.Stdout,
		`digraph {`+"\n"+
			"\t"+`rankdir="LR";`+"\n",
	)
	if err != nil {
		panic(fmt.Errorf("output: %v", err))
	}

	buf := &bytes.Buffer{}

	for _, statement := range allStatements {
		triple := statement.Triple

		sKey := requireResource(buf, "\t", triple.Subject)

		var oKey string

		if oBlankNode, ok := triple.Object.(rdf.BlankNode); ok {
			oKey = requireResource(buf, "\t", oBlankNode)
		} else if oIRI, ok := triple.Object.(rdf.IRI); ok {
			oKey = requireResource(buf, "\t", oIRI)
		} else if oLiteral, ok := triple.Object.(rdf.Literal); ok {
			oKey = fmt.Sprintf("lit%d", literalNodes)
			literalNodes++

			fmt.Fprintf(buf, "\t"+`%s [label=<<table border="1" cellborder="1" cellspacing="0">`+"\n", oKey)
			fmt.Fprintf(buf, "\t\t"+`<tr><td align="left" colspan="2">%s</td></tr>`+"\n", html.EscapeString(oLiteral.LexicalForm))
			fmt.Fprintf(buf, "\t\t"+`<tr><td align="right">Datatype</td><td align="left" href="%s">%s</td></tr>`+"\n", html.EscapeString(string(oLiteral.Datatype)), html.EscapeString(shortenIRI(oLiteral.Datatype)))

			if oLiteral.Tag != nil {
				if tag, ok := oLiteral.Tag.(rdf.LanguageLiteralTag); ok {
					fmt.Fprintf(buf, "\t\t"+`<tr><td align="right">Language</td><td align="left">%s</td></tr>`+"\n", html.EscapeString(tag.Language))
				} else if tag, ok := oLiteral.Tag.(rdf.DirectionalLanguageLiteralTag); ok {
					fmt.Fprintf(buf, "\t\t"+`<tr><td align="right">Language</td><td align="left">%s</td></tr>`+"\n", html.EscapeString(tag.Language))
					fmt.Fprintf(buf, "\t\t"+`<tr><td align="right">Base Direction</td><td align="left">%s</td></tr>`+"\n", html.EscapeString(tag.BaseDirection))
				}
			}

			fmt.Fprintf(buf, "\t"+`</table>>,shape=plain];`+"\n")
		} else {
			panic(fmt.Errorf("invalid object type: %T", triple.Object))
		}

		pIRI := triple.Predicate.(rdf.IRI)

		fmt.Fprintf(buf, "\t"+`%s -> %s [href=%q,label="%s"]`+"\n", sKey, oKey, pIRI, shortenIRI(pIRI))

		_, err = buf.WriteTo(os.Stdout)
		if err != nil {
			panic(fmt.Errorf("output: %v", err))
		}

		buf.Reset()
	}

	_, err = fmt.Fprintf(os.Stdout, `}`+"\n")
	if err != nil {
		panic(fmt.Errorf("output: %v", err))
	}
}
