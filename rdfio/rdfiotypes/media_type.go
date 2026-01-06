package rdfiotypes

import "github.com/dpb587/rdfkit-go/encoding"

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
