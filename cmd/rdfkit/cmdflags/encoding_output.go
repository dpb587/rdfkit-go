package cmdflags

import (
	"context"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
	"github.com/spf13/pflag"
)

type EncodingOutput struct {
	ResourceName   string
	ResourceParams []string

	EncodingName    string
	EncodingBaseIRI string
	EncodingParams  []string
}

func (f *EncodingOutput) Bind(fs *pflag.FlagSet, base, shorthand string) {
	f.BindResource(fs, base, shorthand)
	f.BindEncoding(fs, base, true)
}

func (f *EncodingOutput) BindResource(fs *pflag.FlagSet, base, shorthand string) {
	fs.StringVarP(&f.ResourceName, base, shorthand, f.ResourceName, "path or IRI for writing (default stdout)")
	fs.StringArrayVar(&f.ResourceParams, base+"-param-io", f.ResourceParams, "extra write configuration parameters (syntax \"KEY[=VALUE]\")")
}

func (f *EncodingOutput) BindEncoding(fs *pflag.FlagSet, base string, includeType bool) {
	if includeType {
		fs.StringVar(&f.EncodingName, base+"-type", f.EncodingName, "name or alias for the encoder (default detect or nquads)")
	}

	fs.StringVar(&f.EncodingBaseIRI, base+"-base", f.EncodingBaseIRI, "override the base IRI of the resource")
	fs.StringArrayVar(&f.EncodingParams, base+"-param", f.EncodingParams, "extra encode configuration parameters (syntax \"KEY[=VALUE]\")")
}

type EncodingOutputOpenOptions struct {
	WriterTee           io.Writer
	EncoderPatcher      rdfiotypes.GenericOptionsPatcherFunc
	EncoderDecoderPipe  *rdfiotypes.DecoderHandle
	EncoderFallbackType encoding.ContentTypeIdentifier
}

func (f EncodingOutput) Open(ctx context.Context, r rdfiotypes.Registry, opts *EncodingOutputOpenOptions) (*rdfiotypes.EncoderHandle, error) {
	if opts == nil {
		opts = &EncodingOutputOpenOptions{}
	}

	return r.OpenEncoder(
		ctx,
		rdfiotypes.WriterOptions{
			Name:   f.ResourceName,
			Params: f.ResourceParams,
			Tee:    opts.WriterTee,
		},
		rdfiotypes.EncoderOptions{
			Type:        f.EncodingName,
			BaseIRI:     rdf.IRI(f.EncodingBaseIRI),
			Params:      f.EncodingParams,
			Patcher:     opts.EncoderPatcher,
			DecoderPipe: opts.EncoderDecoderPipe,
		},
		rdfiotypes.EncoderOptionsBuilderFunc(func(r rdfiotypes.Registry, ww rdfiotypes.Writer, ropts *rdfiotypes.EncoderOptions) error {
			cti, ok := r.ResolveEncoderType(ww, ropts.Type)
			if ok {
				ropts.Type = string(cti)

				return nil
			}

			ropts.Type = string(opts.EncoderFallbackType)

			return nil
		}),
	)
}
