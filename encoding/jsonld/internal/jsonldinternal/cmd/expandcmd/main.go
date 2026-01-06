package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

func main() {
	parsed, err := inspectjson.Parse(os.Stdin, inspectjson.TokenizerConfig{}.
		SetSourceOffsets(true),
	)
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

	// e := json.NewEncoder(os.Stdout)
	// e.SetIndent("", "  ")

	// err = e.Encode(expanded.AsBuiltin())
	// if err != nil {
	// 	panic(err)
	// }

	dump(os.Stdout, "", expanded)
}

func dump(w io.Writer, indent string, v any) {
	switch vT := v.(type) {
	case *jsonldinternal.ExpandedArray:
		w.Write([]byte("ExpandedArray{\n"))

		indent1 := indent + "  "

		w.Write([]byte(indent1))
		w.Write([]byte("Values: []ExpandedValue{"))

		indent2 := indent1 + "  "

		for valueIdx, value := range vT.Values {
			if valueIdx > 0 {
				w.Write([]byte(","))
			}

			w.Write([]byte("\n"))
			w.Write([]byte(indent2))

			dump(w, indent2, value)
		}

		if len(vT.Values) > 0 {
			w.Write([]byte("\n"))
			w.Write([]byte(indent1))
		}

		w.Write([]byte("},\n"))
		w.Write([]byte(indent))
	case *jsonldinternal.ExpandedObject:
		w.Write([]byte("ExpandedObject{\n"))

		indent1 := indent + "  "

		if vT.SourceOffsets != nil {
			w.Write([]byte(indent1))
			w.Write([]byte("SourceOffsets: `"))
			w.Write([]byte(vT.SourceOffsets.String()))
			w.Write([]byte("`,\n"))
		}

		if vT.PropertySourceOffsets != nil {
			w.Write([]byte(indent1))
			w.Write([]byte("PropertySourceOffsets: `"))
			w.Write([]byte(vT.PropertySourceOffsets.String()))
			w.Write([]byte("`,\n"))
		}

		w.Write([]byte(indent1))
		w.Write([]byte("Members: map[string]ExpandedValue{"))

		indent2 := indent1 + "  "

		for memberKey, memberValue := range vT.Members {
			w.Write([]byte("\n"))
			w.Write([]byte(indent2))
			w.Write([]byte("\""))
			w.Write([]byte(memberKey))
			w.Write([]byte("\": "))

			dump(w, indent2, memberValue)
		}

		if len(vT.Members) > 0 {
			w.Write([]byte("\n"))
			w.Write([]byte(indent1))
		}

		w.Write([]byte("},\n"))
		w.Write([]byte(indent))
	case *jsonldinternal.ExpandedScalarPrimitive:
		w.Write([]byte("ExpandedScalarPrimitive{\n"))

		indent1 := indent + "  "

		if vT.PropertySourceOffsets != nil {
			w.Write([]byte(indent1))
			w.Write([]byte("PropertySourceOffsets: `"))
			w.Write([]byte(vT.PropertySourceOffsets.String()))
			w.Write([]byte("`,\n"))
		}

		w.Write([]byte(indent1))
		w.Write([]byte("Value: "))

		buf, err := json.Marshal(vT.Value.AsBuiltin())
		if err != nil {
			panic(err)
		}

		w.Write(buf)
		w.Write([]byte(",\n"))

		w.Write([]byte(indent))
		w.Write([]byte("}"))
	default:
		panic(fmt.Errorf("unknown type: %T", vT))
	}
}
