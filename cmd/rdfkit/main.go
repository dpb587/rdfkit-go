package main

import (
	"net/http"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/canonicalizecmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/inspectdecodercmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/irigencmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/pipecmd"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/x/encodingref"
	"github.com/dpb587/rdfkit-go/x/encodingref/encodingdefaults"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "rdfkit",
	}

	resourceManager := encodingref.NewResourceManager(
		encodingref.NewWebResourceManager(http.DefaultClient),
		encodingref.NewFileResourceManager(),
	)

	encodingRegistry := encodingdefaults.NewRegistry(encodingdefaults.RegistryOptions{
		DocumentLoaderJSONLD: jsonldtype.NewCachingDocumentLoader(
			jsonldtype.NewDefaultDocumentLoader(
				http.DefaultClient,
			),
		),
	})

	cmd.AddCommand(
		pipecmd.New(resourceManager, encodingRegistry),
		canonicalizecmd.New(resourceManager, encodingRegistry),
		irigencmd.New(resourceManager, encodingRegistry),
		inspectdecodercmd.New(resourceManager, encodingRegistry),
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
