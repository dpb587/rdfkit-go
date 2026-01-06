package pipecmd

import (
	"context"
	"fmt"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/spf13/cobra"
)

func New(app *cmdutil.App) *cobra.Command {
	fIn := &cmdflags.EncodingInput{}
	fOut := &cmdflags.EncodingOutput{}

	cmd := &cobra.Command{
		Use:   "pipe",
		Short: "Decode and re-encode using supported encoding formats",
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

			bfOut, err := fOut.Open(ctx, app.Registry, &cmdflags.EncodingOutputOpenOptions{
				EncoderDecoderPipe:  bfIn,
				EncoderFallbackType: nquadscontent.TypeIdentifier,
			})
			if err != nil {
				return fmt.Errorf("output: %v", err)
			}

			defer bfOut.Close()

			decoderQuads := bfIn.GetQuadsDecoder()
			encoderQuads := bfOut.GetQuadsEncoder()

			writeStatementFunc := func(ctx context.Context, iter rdf.QuadIterator) error {
				return encoderQuads.AddQuad(ctx, iter.Quad())
			}

			if statementEncoder, ok := encoderQuads.(interface {
				AddQuadStatement(context.Context, encodingtest.QuadStatement) error
			}); ok {
				if statementDecoder, ok := decoderQuads.(encoding.StatementTextOffsetsProvider); ok {
					writeStatementFunc = func(ctx context.Context, iter rdf.QuadIterator) error {
						return statementEncoder.AddQuadStatement(ctx, encodingtest.QuadStatement{
							Quad:        iter.Quad(),
							TextOffsets: statementDecoder.StatementTextOffsets(),
						})
					}
				}
			}

			for decoderQuads.Next() {
				err := writeStatementFunc(ctx, decoderQuads)
				if err != nil {
					return fmt.Errorf("write: %v", err)
				}
			}

			if err := decoderQuads.Err(); err != nil {
				return fmt.Errorf("decode[%s]: %v", bfIn.Decoder.GetContentTypeIdentifier(), err)
			}

			return nil
		},
	}

	f := cmd.Flags()
	fIn.Bind(f, "in", "i")
	fOut.Bind(f, "out", "o")

	cmd.SetHelpFunc(cmdutil.RegistryHelpFunc(app))

	return cmd
}
