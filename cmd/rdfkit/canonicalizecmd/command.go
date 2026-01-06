package canonicalizecmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dpb587/rdfkit-go/cmd/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdfcanon"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	fIn := &cmdflags.EncodingInput{
		Path:           "-",
		FallbackOpener: cmdflags.WebRemoteOpener,
		DocumentLoaderJSONLD: jsonldtype.NewCachingDocumentLoader(
			jsonldtype.NewDefaultDocumentLoader(
				http.DefaultClient,
			),
		),
	}

	fOut := &cmdflags.EncodingOutput{
		Path: "-",
		Type: "nquads",
	}

	cmd := &cobra.Command{
		Use: "canonicalize",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			bfIn, err := fIn.Open()
			if err != nil {
				return fmt.Errorf("input: %v", err)
			}

			defer bfIn.Close()

			rdfc, err := rdfcanon.Canonicalize(ctx, bfIn.Decoder)
			if err != nil {
				return fmt.Errorf("canonicalize: %v", err)
			}

			_, err = rdfc.WriteTo(os.Stdout)
			if err != nil {
				return fmt.Errorf("write: %v", err)
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&fIn.Path, "in", "i", fIn.Path, "")
	f.StringVar(&fIn.Type, "in-type", fIn.Type, "")
	f.StringVar(&fIn.DefaultBase, "in-default-base", fIn.DefaultBase, "")
	f.StringVarP(&fOut.Path, "out", "o", fOut.Path, "")

	return cmd
}
