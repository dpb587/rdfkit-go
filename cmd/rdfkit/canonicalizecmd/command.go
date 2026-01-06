package canonicalizecmd

import (
	"fmt"
	"os"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdfcanon"
	"github.com/spf13/cobra"
)

func New(app *cmdutil.App) *cobra.Command {
	fIn := &cmdflags.EncodingInput{}

	cmd := &cobra.Command{
		Use:   "canonicalize",
		Short: "Convert a dataset into canonical blank nodes and ordering",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			bfIn, err := fIn.Open(ctx, app.Registry, &cmdflags.EncodingInputOpenOptions{
				DecoderFallbackType: trigcontent.TypeIdentifier,
			})
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

	return cmd
}
