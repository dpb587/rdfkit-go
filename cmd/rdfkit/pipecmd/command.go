package pipecmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dpb587/rdfkit-go/cmd/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
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
		Use: "pipe",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			bfIn, err := fIn.Open()
			if err != nil {
				return fmt.Errorf("input: %v", err)
			}

			defer bfIn.Close()

			bfOut, err := fOut.NewStatementWriter()
			if err != nil {
				return fmt.Errorf("output: %v", err)
			}

			defer bfOut.Close()

			writeStatementFunc := func(ctx context.Context, iter rdf.QuadIterator) error {
				return bfOut.Encoder.AddQuad(ctx, iter.Quad())
			}

			var seqCloser func() error

			if descriptionsEncoder, ok := bfOut.Encoder.(rdfdescription.DatasetResourceWriter); ok {
				descriptions := rdfdescription.NewDatasetResourceListBuilder()

				writeStatementFunc = func(_ context.Context, iter rdf.QuadIterator) error {
					descriptions.Add(iter.Quad())

					return nil
				}

				seqCloser = func() error {
					return descriptions.AddToDataset(ctx, descriptionsEncoder, true)
				}
			} else if statementEncoder, ok := bfOut.Encoder.(interface {
				AddQuadStatement(context.Context, encodingtest.QuadStatement) error
			}); ok {
				if statementDecoder, ok := bfIn.Decoder.(encoding.StatementTextOffsetsProvider); ok {
					writeStatementFunc = func(ctx context.Context, iter rdf.QuadIterator) error {
						return statementEncoder.AddQuadStatement(ctx, encodingtest.QuadStatement{
							Quad:        iter.Quad(),
							TextOffsets: statementDecoder.StatementTextOffsets(),
						})
					}
				}
			}

			for bfIn.Decoder.Next() {
				err := writeStatementFunc(ctx, bfIn.Decoder)
				if err != nil {
					return fmt.Errorf("write: %v", err)
				}
			}

			if err := bfIn.Decoder.Err(); err != nil {
				return fmt.Errorf("read: %s: %v", bfIn.Format, err)
			}

			if seqCloser != nil {
				if err := seqCloser(); err != nil {
					return fmt.Errorf("write: %v", err)
				}
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&fIn.Path, "in", "i", fIn.Path, "")
	f.StringVar(&fIn.Type, "in-type", fIn.Type, "")
	f.StringVar(&fIn.DefaultBase, "in-default-base", fIn.DefaultBase, "")
	f.BoolVar(&fIn.SkipTextOffsets, "in-skip-text-offsets", fIn.SkipTextOffsets, "")
	f.StringVarP(&fOut.Path, "out", "o", fOut.Path, "")
	f.StringVar(&fOut.Type, "out-type", fOut.Type, "")

	return cmd
}
