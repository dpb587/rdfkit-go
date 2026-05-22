package exportgoiricmd

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/spf13/cobra"
)

func New(app *cmdutil.App) *cobra.Command {
	var pkgName = ""
	var base string
	var prefix string
	var outPath string

	fIn := &cmdflags.EncodingInput{}

	cmd := &cobra.Command{
		Use:   "export-go-iri",
		Short: "Generate a Go file of IRI constants from an ontology",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if len(outPath) == 0 {
				return errors.New("output path required")
			}

			hasher := sha256.New()

			bfIn, err := fIn.Open(ctx, app.Registry, &cmdflags.EncodingInputOpenOptions{
				ReaderTee:           hasher,
				DecoderFallbackType: trigcontent.TypeIdentifier,
			})
			if err != nil {
				return fmt.Errorf("input: %v", err)
			}

			defer bfIn.Close()

			if len(pkgName) == 0 {
				workdir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("getwd: %v", err)
				}

				if len(outPath) > 0 {
					pkgName = filepath.Base(filepath.Dir(filepath.Join(workdir, outPath)))
				} else {
					pkgName = filepath.Base(workdir)
				}

				if pkgName == "." {
					pkgName = "iri"
				}
			}

			b, err := NewBuilderFromQuadsDecoder(bfIn.GetQuadsDecoder())
			if err != nil {
				return fmt.Errorf("decode[%s]: %v", bfIn.Decoder.GetContentTypeIdentifier(), err)
			}

			if len(base) == 0 {
				detectedBase, ok := b.DetectBase()
				if !ok {
					return errors.New("unable to detect base IRI")
				}

				base = detectedBase
			}

			b = b.FilterBase(rdf.IRI(base))

			if len(b.statementsBySubject) == 0 {
				return errors.New("no statements found")
			}

			var out io.Writer
			if outPath == "-" {
				out = os.Stdout
			} else {
				fh, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					return fmt.Errorf("open output: %v", err)
				}

				defer fh.Close()

				out = fh
			}

			if err := b.ExportSource(out, ExportSourceOptions{
				PackageName: pkgName,
				Prefix:      prefix,
				SourceHash:  hasher.Sum(nil),
				SourceIRI:   string(bfIn.Reader.GetIRI()),
			}); err != nil {
				return fmt.Errorf("export: %v", err)
			}

			return nil
		},
	}

	f := cmd.Flags()
	fIn.Bind(f, "in", "i")
	f.StringVar(&base, "base", "", "vocab base IRI")
	f.StringVar(&prefix, "prefix", "", "optional prefix for tokens")
	f.StringVarP(&outPath, "out", "o", outPath, "go output file")

	return cmd
}
