package cmdflags

import (
	"context"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
	"github.com/spf13/pflag"
)

type EncodingInput struct {
	ResourceName   string
	ResourceParams []string

	EncodingName    string
	EncodingBaseIRI string
	EncodingParams  []string
}

func (f *EncodingInput) Bind(fs *pflag.FlagSet, base, shorthand string) {
	f.BindResource(fs, base, shorthand)
	f.BindEncoding(fs, base, true)
}

func (f *EncodingInput) BindResource(fs *pflag.FlagSet, base, shorthand string) {
	fs.StringVarP(&f.ResourceName, base, shorthand, f.ResourceName, "path or IRI for reading (default stdin)")
	fs.StringArrayVar(&f.ResourceParams, base+"-param-io", f.ResourceParams, "extra read configuration parameters (syntax \"KEY[=VALUE]\")")
}

func (f *EncodingInput) BindEncoding(fs *pflag.FlagSet, base string, includeType bool) {
	if includeType {
		fs.StringVar(&f.EncodingName, base+"-type", f.EncodingName, "name or alias for the decoder (default detect)")
	}

	fs.StringVar(&f.EncodingBaseIRI, base+"-base", f.EncodingBaseIRI, "override the base IRI of the resource")
	fs.StringArrayVar(&f.EncodingParams, base+"-param", f.EncodingParams, "extra decode configuration parameters (syntax \"KEY[=VALUE]\")")
}

type EncodingInputOpenOptions struct {
	ReaderTee           io.Writer
	DecoderPatcher      rdfiotypes.GenericOptionsPatcherFunc
	DecoderFallbackType encoding.ContentTypeIdentifier
}

func (f EncodingInput) Open(ctx context.Context, r rdfiotypes.Registry, opts *EncodingInputOpenOptions) (*rdfiotypes.DecoderHandle, error) {
	if opts == nil {
		opts = &EncodingInputOpenOptions{}
	}

	return r.OpenDecoder(
		ctx,
		rdfiotypes.ReaderOptions{
			Name:   f.ResourceName,
			Params: f.ResourceParams,
			Tee:    opts.ReaderTee,
		},
		rdfiotypes.DecoderOptions{
			Type:    f.EncodingName,
			BaseIRI: rdf.IRI(f.EncodingBaseIRI),
			Params:  f.EncodingParams,
			Patcher: opts.DecoderPatcher,
		},
		rdfiotypes.DecoderOptionsBuilderFunc(func(r rdfiotypes.Registry, rr rdfiotypes.Reader, ropts *rdfiotypes.DecoderOptions) error {
			cti, ok := r.ResolveDecoderType(rr, ropts.Type)
			if ok {
				ropts.Type = string(cti)

				return nil
			}

			ropts.Type = string(opts.DecoderFallbackType)

			return nil
		}),
	)
}
