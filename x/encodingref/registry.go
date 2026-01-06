package encodingref

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"path/filepath"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

var ErrEncodingNotSupported = errors.New("encoding not supported")

//

type DecoderOptions struct {
	IRI    rdf.IRI
	Flags  []string
	Tee    io.Writer
	Hasher hash.Hash
}

type EncoderOptions struct {
	IRI    rdf.IRI
	Flags  []string
	Tee    io.Writer
	Hasher hash.Hash
}

type DecoderHandle struct {
	Reader  ResourceReader
	Decoder encoding.Decoder

	DecodedBase           []string
	DecodedPrefixMappings iriutil.PrefixMappingList
}

func (b *DecoderHandle) Close() error {
	if err := b.Decoder.Close(); err != nil {
		return err
	}

	if err := b.Reader.Close(); err != nil {
		return err
	}

	return nil
}

func (dh *DecoderHandle) GetQuadsDecoder() encoding.QuadsDecoder {
	switch d := dh.Decoder.(type) {
	case encoding.QuadsDecoder:
		return d
	case encoding.TriplesDecoder:
		return encodingutil.NewTripleAsQuadDecoder(d, nil)
	}

	panic(fmt.Errorf("unexpected decoder type: %T", dh.Decoder))
}

type EncoderHandle struct {
	Writer  ResourceWriter
	Encoder encoding.Encoder
}

func (b *EncoderHandle) Close() error {
	if err := b.Encoder.Close(); err != nil {
		return err
	}

	if err := b.Writer.Close(); err != nil {
		return err
	}

	return nil
}

func (dh *EncoderHandle) GetQuadsEncoder() encoding.QuadsEncoder {
	switch d := dh.Encoder.(type) {
	case encoding.QuadsEncoder:
		return d
	case encoding.TriplesEncoder:
		return encodingutil.QuadAsTripleEncoder{
			TriplesEncoder: d,
		}
	}

	panic(fmt.Errorf("unexpected encoder type: %T", dh.Encoder))
}

//

type RegistryEncoding interface {
	NewDecoder(cti encoding.ContentTypeIdentifier, rr ResourceReader, opts DecoderOptions) (*DecoderHandle, error)
	NewEncoder(cti encoding.ContentTypeIdentifier, rw ResourceWriter, opts EncoderOptions) (*EncoderHandle, error)
}

//

type RegistryOptions struct {
	Aliases             map[string]encoding.ContentTypeIdentifier
	MediaTypes          map[string]encoding.ContentTypeIdentifier
	FileExts            map[string]encoding.ContentTypeIdentifier
	MagicBytesResolvers []MagicBytesResolver
	Encodings           map[encoding.ContentTypeIdentifier]RegistryEncoding
}

//

type MediaTypeResolver interface {
	ResolveMediaType(mt encoding.ContentMediaType) (encoding.ContentTypeIdentifier, bool)
}

//

type MediaTypeResolverFunc func(mt encoding.ContentMediaType) (encoding.ContentTypeIdentifier, bool)

func (f MediaTypeResolverFunc) ResolveMediaType(mt encoding.ContentMediaType) (encoding.ContentTypeIdentifier, bool) {
	return f(mt)
}

//

type MagicBytesResolver interface {
	ResolveMagicBytes(mt []byte) (encoding.ContentTypeIdentifier, bool)
}

//

type MagicBytesResolverFunc func(buf []byte) (encoding.ContentTypeIdentifier, bool)

func (f MagicBytesResolverFunc) ResolveMagicBytes(buf []byte) (encoding.ContentTypeIdentifier, bool) {
	return f(buf)
}

//

type Registry struct {
	aliases             map[string]encoding.ContentTypeIdentifier
	mediaTypes          map[string]encoding.ContentTypeIdentifier
	fileExts            map[string]encoding.ContentTypeIdentifier
	magicBytesResolvers []MagicBytesResolver
	encodings           map[encoding.ContentTypeIdentifier]RegistryEncoding
}

func NewRegistry(opts RegistryOptions) Registry {
	r := Registry{
		aliases:             opts.Aliases,
		mediaTypes:          opts.MediaTypes,
		fileExts:            opts.FileExts,
		magicBytesResolvers: opts.MagicBytesResolvers,
		encodings:           opts.Encodings,
	}

	return r
}

func (r Registry) ResolveName(name string) (encoding.ContentTypeIdentifier, bool) {
	if cti, ok := r.aliases[name]; ok {
		return cti, true
	}

	if _, ok := r.encodings[encoding.ContentTypeIdentifier(name)]; ok {
		return encoding.ContentTypeIdentifier(name), true
	}

	return "", false
}

func (r Registry) ResolveReader(rr ResourceReader) (encoding.ContentTypeIdentifier, bool) {
	if mt, ok := rr.GetMediaType(); ok {
		if cti, ok := r.mediaTypes[strings.ToLower(mt.Type+"/"+mt.Subtype)]; ok {
			return cti, true
		}
	}

	if mb, ok := rr.GetMagicBytes(); ok {
		for _, resolver := range r.magicBytesResolvers {
			if cti, ok := resolver.ResolveMagicBytes(mb); ok {
				return cti, true
			}
		}
	}

	if fileName, ok := rr.GetFileName(); ok {
		if cti, ok := r.fileExts[filepath.Ext(fileName)]; ok {
			return cti, true
		}
	}

	return "", false
}

func (r Registry) ResolveWriter(ww ResourceWriter) (encoding.ContentTypeIdentifier, bool) {
	if fileName, ok := ww.GetFileName(); ok {
		if cti, ok := r.fileExts[filepath.Ext(fileName)]; ok {
			return cti, true
		}
	}

	return "", false
}

func (r Registry) GetEncoding(cti encoding.ContentTypeIdentifier) (RegistryEncoding, bool) {
	encoding, ok := r.encodings[cti]

	return encoding, ok
}

func (r Registry) NewDecoder(cti encoding.ContentTypeIdentifier, rr ResourceReader, opts DecoderOptions) (*DecoderHandle, error) {
	encoding, ok := r.encodings[cti]
	if !ok {
		return nil, ErrEncodingNotSupported
	}

	return encoding.NewDecoder(cti, rr, opts)
}

func (r Registry) NewEncoder(cti encoding.ContentTypeIdentifier, rw ResourceWriter, opts EncoderOptions) (*EncoderHandle, error) {
	encoding, ok := r.encodings[cti]
	if !ok {
		return nil, ErrEncodingNotSupported
	}

	return encoding.NewEncoder(cti, rw, opts)
}
