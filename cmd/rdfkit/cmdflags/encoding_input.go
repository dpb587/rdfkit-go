package cmdflags

import (
	"context"
	"errors"
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/x/encodingref"
	"github.com/spf13/pflag"
)

type EncodingInput struct {
	ResourceName    string
	ResourceOptions []string
	ResourceIRI     string

	EncodingName         string
	EncodingOptions      []string
	EncodingFallbackType encoding.ContentTypeIdentifier
}

func (f *EncodingInput) Bind(fs *pflag.FlagSet, base, shorthand string) {
	fs.StringVarP(&f.ResourceName, base, shorthand, f.ResourceName, "")
	fs.StringVar(&f.ResourceIRI, base+"-base", f.ResourceIRI, "")
	fs.StringArrayVar(&f.ResourceOptions, base+"-io-option", f.ResourceOptions, "")
	fs.StringVar(&f.EncodingName, base+"-type", f.EncodingName, "")
	fs.StringArrayVar(&f.EncodingOptions, base+"-option", f.EncodingOptions, "")
}

func (f EncodingInput) Open(ctx context.Context, resourceManager encodingref.ResourceManager, encodingRegistry encodingref.Registry) (*encodingref.DecoderHandle, error) {
	return f.OpenOptions(ctx, resourceManager, encodingRegistry, encodingref.DecoderOptions{})
}

func (f EncodingInput) OpenOptions(ctx context.Context, resourceManager encodingref.ResourceManager, encodingRegistry encodingref.Registry, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	rr, err := resourceManager.OpenReader(ctx, encodingref.ResourceRef{
		Name:  f.ResourceName,
		Flags: f.ResourceOptions,
	})
	if err != nil {
		return nil, fmt.Errorf("open resource: %v", err)
	}

	var cti encoding.ContentTypeIdentifier
	var ok bool

	if len(f.EncodingName) > 0 {
		cti, ok = encodingRegistry.ResolveName(f.EncodingName)
		if !ok {
			rr.Close()

			return nil, fmt.Errorf("unknown encoding: %s", f.EncodingName)
		}
	}

	if len(cti) == 0 {
		cti, ok = encodingRegistry.ResolveReader(rr)
		if !ok {
			if len(f.EncodingFallbackType) > 0 {
				cti = f.EncodingFallbackType
			} else {
				rr.Close()

				return nil, errors.New("failed to detect encoding")
			}
		}
	}

	if len(f.ResourceIRI) > 0 {
		opts.IRI = rdf.IRI(f.ResourceIRI)
	}

	opts.Flags = append(f.EncodingOptions, opts.Flags...)

	eh, err := encodingRegistry.NewDecoder(cti, rr, opts)
	if err != nil {
		rr.Close()

		return nil, fmt.Errorf("open decoder: %v", err)
	}

	return eh, nil
}
