package pipecmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dpb587/rdfkit-go/cmd/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionio"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
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

	var (
		fGraphAction     string = ""
		fOutDescriptions bool
	)

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

			writeStatementFunc := func(ctx context.Context, s rdfio.Statement) error {
				return bfOut.Encoder.PutTriple(ctx, s.GetTriple())
			}

			if fGraphAction == "drop" {
				if _, ok := bfOut.Encoder.(encoding.DatasetEncoder); !ok {
					bfIn.Decoder = rdfioutil.NewStatementMatcherIterator(bfIn.Decoder, nil)
				}
			}

			var seqCloser func() error

			if _, ok := bfIn.Decoder.(encoding.DatasetDecoder); ok {
				if datasetEncoder, ok := bfOut.Encoder.(encoding.DatasetEncoder); !ok && fGraphAction != "drop" {
					return fmt.Errorf("output must be a dataset writer")
				} else if descriptionsEncoder, ok := bfOut.Encoder.(rdfdescriptionio.DatasetEncoder); ok && fOutDescriptions {
					descriptionsByGraphName := map[rdf.GraphNameValue]*rdfdescription.ResourceListBuilder{}

					writeStatementFunc = func(_ context.Context, s rdfio.Statement) error {
						gn := s.GetGraphName()

						descriptions, ok := descriptionsByGraphName[gn]
						if !ok {
							descriptions = rdfdescription.NewResourceListBuilder()
							descriptionsByGraphName[gn] = descriptions
						}

						descriptions.AddTriple(s.GetTriple())

						return nil
					}

					seqCloser = func() error {
						for gn, descriptions := range descriptionsByGraphName {
							for _, r := range descriptions.GetResources() {
								err := descriptionsEncoder.PutGraphResource(ctx, gn, r)
								if err != nil {
									return err
								}
							}
						}

						return nil
					}
				} else if statementEncoder, ok := bfOut.Encoder.(interface {
					PutStatement(ctx context.Context, s rdfio.Statement) error
				}); ok {
					writeStatementFunc = func(ctx context.Context, s rdfio.Statement) error {
						return statementEncoder.PutStatement(ctx, s)
					}
				} else {
					writeStatementFunc = func(ctx context.Context, s rdfio.Statement) error {
						return datasetEncoder.PutGraphTriple(ctx, s.GetGraphName(), s.GetTriple())
					}
				}
			} else if descriptionsEncoder, ok := bfOut.Encoder.(rdfdescriptionio.GraphEncoder); ok && fOutDescriptions {
				descriptions := rdfdescription.NewResourceListBuilder()

				writeStatementFunc = func(_ context.Context, s rdfio.Statement) error {
					descriptions.AddTriple(s.GetTriple())

					return nil
				}

				seqCloser = func() error {
					for _, r := range descriptions.GetResources() {
						err := descriptionsEncoder.PutResource(ctx, r)
						if err != nil {
							return err
						}
					}

					return nil
				}
			} else if statementEncoder, ok := bfOut.Encoder.(interface {
				PutStatement(ctx context.Context, s rdfio.Statement) error
			}); ok {
				writeStatementFunc = func(ctx context.Context, s rdfio.Statement) error {
					return statementEncoder.PutStatement(ctx, s)
				}
			}

			for bfIn.Decoder.Next() {
				err := writeStatementFunc(ctx, bfIn.Decoder.GetStatement())
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
	f.StringVar(&fGraphAction, "graph-action", fGraphAction, "")
	f.StringVarP(&fOut.Path, "out", "o", fOut.Path, "")
	f.StringVar(&fOut.Type, "out-type", fOut.Type, "")
	f.BoolVar(&fOutDescriptions, "out-descriptions", fOutDescriptions, "")

	return cmd
}
