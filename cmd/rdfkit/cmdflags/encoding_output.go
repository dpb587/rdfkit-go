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

type EncodingOutput struct {
	ResourceName    string
	ResourceOptions []string
	ResourceIRI     string

	EncodingName         string
	EncodingOptions      []string
	EncodingFallbackType encoding.ContentTypeIdentifier
}

func (f *EncodingOutput) Bind(fs *pflag.FlagSet, base, shorthand string) {
	fs.StringVarP(&f.ResourceName, base, shorthand, f.ResourceName, "")
	fs.StringVar(&f.ResourceIRI, base+"-base", f.ResourceIRI, "")
	fs.StringArrayVar(&f.ResourceOptions, base+"-io-option", f.ResourceOptions, "")
	fs.StringVar(&f.EncodingName, base+"-type", f.EncodingName, "")
	fs.StringArrayVar(&f.EncodingOptions, base+"-option", f.EncodingOptions, "")
}

func (f EncodingOutput) Open(ctx context.Context, resourceManager encodingref.ResourceManager, encodingRegistry encodingref.Registry) (*encodingref.EncoderHandle, error) {
	return f.OpenOptions(ctx, resourceManager, encodingRegistry, encodingref.EncoderOptions{})
}

func (f EncodingOutput) OpenOptions(ctx context.Context, resourceManager encodingref.ResourceManager, encodingRegistry encodingref.Registry, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	ww, err := resourceManager.OpenWriter(ctx, encodingref.ResourceRef{
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
			ww.Close()

			return nil, fmt.Errorf("unknown encoding: %s", f.EncodingName)
		}
	}

	if len(cti) == 0 {
		cti, ok = encodingRegistry.ResolveWriter(ww)
		if !ok {
			if len(f.EncodingFallbackType) > 0 {
				cti = f.EncodingFallbackType
			} else {
				ww.Close()

				return nil, errors.New("failed to detect encoding")
			}
		}
	}

	if len(f.ResourceIRI) > 0 {
		opts.IRI = rdf.IRI(f.ResourceIRI)
	}

	opts.Flags = append(f.EncodingOptions, opts.Flags...)

	eh, err := encodingRegistry.NewEncoder(cti, ww, opts)
	if err != nil {
		ww.Close()

		return nil, fmt.Errorf("open decoder: %v", err)
	}

	return eh, nil
}
