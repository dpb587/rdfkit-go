package pipecmd

import (
	"context"
	"fmt"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdf"
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
		Use: "pipe",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			bfIn, err := fIn.Open(ctx, resourceManager, encodingRegistry)
			if err != nil {
				return fmt.Errorf("input: %v", err)
			}

			defer bfIn.Close()

			bfOut, err := fOut.Open(ctx, resourceManager, encodingRegistry)
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

	return cmd
}
