package rdfiotypes

import (
	"context"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding"
)

//

type Registry struct {
	ResourceManagers []ResourceManager
	EncoderManagers  map[encoding.ContentTypeIdentifier]EncoderManager
	DecoderManagers  map[encoding.ContentTypeIdentifier]DecoderManager

	Aliases    map[string]encoding.ContentTypeIdentifier
	MediaTypes map[string]encoding.ContentTypeIdentifier

	// FileExts contains a map from file extensions to content type identifiers. File extensions must be lowercase,
	// should include a leading dot, and are matched as suffixes against file names.
	FileExts            map[string]encoding.ContentTypeIdentifier
	MagicBytesResolvers []MagicBytesResolver
}

type RegistryMapper interface {
	MapRegistry(r Registry) Registry
}

type RegistryMapperFunc func(r Registry) Registry

func (f RegistryMapperFunc) MapRegistry(r Registry) Registry {
	return f(r)
}

func (r Registry) Clone(m ...RegistryMapper) Registry {
	next := Registry{
		ResourceManagers: slices.Clone(r.ResourceManagers),
		EncoderManagers:  maps.Clone(r.EncoderManagers),
		DecoderManagers:  maps.Clone(r.DecoderManagers),

		Aliases:             maps.Clone(r.Aliases),
		MediaTypes:          maps.Clone(r.MediaTypes),
		FileExts:            maps.Clone(r.FileExts),
		MagicBytesResolvers: slices.Clone(r.MagicBytesResolvers),
	}

	for _, fn := range m {
		next = fn.MapRegistry(next)
	}

	return next
}

func (r Registry) ResolveDecoderType(rr Reader, t string) (encoding.ContentTypeIdentifier, bool) {
	if len(t) > 0 {
		if cti, ok := r.Aliases[t]; ok {
			return cti, true
		} else if _, ok := r.DecoderManagers[encoding.ContentTypeIdentifier(t)]; ok {
			return encoding.ContentTypeIdentifier(t), true
		}
	}

	if mt, ok := rr.GetMediaType(); ok {
		if cti, ok := r.MediaTypes[strings.ToLower(mt.Type+"/"+mt.Subtype)]; ok {
			return cti, true
		}
	}

	if mb, ok := rr.GetMagicBytes(); ok {
		for _, resolver := range r.MagicBytesResolvers {
			if cti, ok := resolver.ResolveMagicBytes(mb); ok {
				return cti, true
			}
		}
	}

	if fileName, ok := rr.GetFileName(); ok {
		fileNameLower := strings.ToLower(fileName)

		for fileExt, cti := range r.FileExts {
			if strings.HasSuffix(fileNameLower, fileExt) {
				return cti, true
			}
		}
	}

	return "", false
}

func (r Registry) ResolveEncoderType(ww Writer, t string) (encoding.ContentTypeIdentifier, bool) {
	if len(t) > 0 {
		if cti, ok := r.Aliases[t]; ok {
			return cti, true
		} else if _, ok := r.EncoderManagers[encoding.ContentTypeIdentifier(t)]; ok {
			return encoding.ContentTypeIdentifier(t), true
		}
	}

	if fileName, ok := ww.GetFileName(); ok {
		if cti, ok := r.FileExts[filepath.Ext(fileName)]; ok {
			return cti, true
		}
	}

	return "", false
}

func (r Registry) OpenReader(ctx context.Context, opts ReaderOptions) (Reader, error) {
	for _, rm := range r.ResourceManagers {
		rr, err := rm.NewReader(ctx, opts)
		if err == ErrResourceNotSupported {
			continue
		} else if err != nil {
			return nil, err
		}

		return rr, nil
	}

	return nil, ErrResourceNotSupported
}

func (r Registry) OpenWriter(ctx context.Context, opts WriterOptions) (Writer, error) {
	for _, rm := range r.ResourceManagers {
		ww, err := rm.NewWriter(ctx, opts)
		if err == ErrResourceNotSupported {
			continue
		} else if err != nil {
			return nil, err
		}

		return ww, nil
	}

	return nil, ErrResourceNotSupported
}

func (r Registry) NewDecoder(rr Reader, opts ...DecoderOptionsBuilder) (*DecoderHandle, error) {
	applied := DecoderOptions{}

	for _, opt := range opts {
		if err := opt.ApplyOptions(r, rr, &applied); err != nil {
			return nil, fmt.Errorf("decoder options: %v", err)
		}
	}

	if len(applied.BaseIRI) == 0 {
		applied.BaseIRI = rr.GetIRI()
	}

	cti, ok := r.ResolveDecoderType(rr, applied.Type)
	if !ok {
		return nil, ErrUnknownEncoding
	}

	decoderManager, ok := r.DecoderManagers[cti]
	if !ok {
		return nil, ErrUnknownEncoding
	}

	return decoderManager.NewDecoder(rr, applied)
}

func (r Registry) NewEncoder(ww Writer, opts ...EncoderOptionsBuilder) (*EncoderHandle, error) {
	applied := EncoderOptions{}

	for _, opt := range opts {
		if err := opt.ApplyOptions(r, ww, &applied); err != nil {
			return nil, fmt.Errorf("decoder options: %v", err)
		}
	}

	if len(applied.BaseIRI) == 0 {
		applied.BaseIRI = ww.GetIRI()
	}

	cti, ok := r.ResolveEncoderType(ww, applied.Type)
	if !ok {
		return nil, ErrUnknownEncoding
	}

	encoderManager, ok := r.EncoderManagers[cti]
	if !ok {
		return nil, ErrUnknownEncoding
	}

	return encoderManager.NewEncoder(ww, applied)
}

type DecoderOptionsBuilder interface {
	ApplyOptions(r Registry, rr Reader, opts *DecoderOptions) error
}

type DecoderOptionsBuilderFunc func(r Registry, rr Reader, opts *DecoderOptions) error

func (f DecoderOptionsBuilderFunc) ApplyOptions(r Registry, rr Reader, opts *DecoderOptions) error {
	return f(r, rr, opts)
}

func (r Registry) OpenDecoder(ctx context.Context, ropts ReaderOptions, dopts ...DecoderOptionsBuilder) (*DecoderHandle, error) {
	rr, err := r.OpenReader(ctx, ropts)
	if err != nil {
		return nil, fmt.Errorf("reader: %v", err)
	}

	d, err := r.NewDecoder(rr, dopts...)
	if err != nil {
		rr.Close()

		return nil, fmt.Errorf("decoder: %v", err)
	}

	return d, nil
}

type EncoderOptionsBuilder interface {
	ApplyOptions(r Registry, ww Writer, opts *EncoderOptions) error
}

type EncoderOptionsBuilderFunc func(r Registry, ww Writer, opts *EncoderOptions) error

func (f EncoderOptionsBuilderFunc) ApplyOptions(r Registry, ww Writer, opts *EncoderOptions) error {
	return f(r, ww, opts)
}

func (r Registry) OpenEncoder(ctx context.Context, wopts WriterOptions, eopts ...EncoderOptionsBuilder) (*EncoderHandle, error) {
	ww, err := r.OpenWriter(ctx, wopts)
	if err != nil {
		return nil, fmt.Errorf("writer: %v", err)
	}

	d, err := r.NewEncoder(ww, eopts...)
	if err != nil {
		ww.Close()

		return nil, fmt.Errorf("encoder: %v", err)
	}

	return d, nil
}
