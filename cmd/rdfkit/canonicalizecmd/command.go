package canonicalizecmd

import (
	"fmt"
	"os"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdfcanon"
	"github.com/dpb587/rdfkit-go/x/encodingref"
	"github.com/spf13/cobra"
)

func New(resourceManager encodingref.ResourceManager, encodingRegistry encodingref.Registry) *cobra.Command {
	fIn := &cmdflags.EncodingInput{
		ResourceName:         "-",
		EncodingFallbackType: trigcontent.TypeIdentifier,
	}

	fOut := &cmdflags.EncodingOutput{
		ResourceName:         "-",
		EncodingFallbackType: nquadscontent.TypeIdentifier,
	}

	cmd := &cobra.Command{
		Use: "canonicalize",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			bfIn, err := fIn.Open(ctx, resourceManager, encodingRegistry)
			if err != nil {
				return fmt.Errorf("input: %v", err)
			}

			defer bfIn.Close()

			rdfc, err := rdfcanon.Canonicalize(ctx, bfIn.GetQuadsDecoder())
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
	fIn.Bind(f, "in", "i")
	fOut.Bind(f, "out", "o")

	return cmd
}
