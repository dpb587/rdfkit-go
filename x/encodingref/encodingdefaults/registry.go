package encodingdefaults

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/html/htmlcontent"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldcontent"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/ntriplescontent"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson/rdfjsoncontent"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml/rdfxmlcontent"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/encoding/turtle/turtlecontent"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type RegistryOptions struct {
	DocumentLoaderJSONLD jsonldtype.DocumentLoader
}

func NewRegistry(opts RegistryOptions) encodingref.Registry {
	devEncoding := encodingDev{}

	return encodingref.NewRegistry(encodingref.RegistryOptions{
		Aliases: map[string]encoding.ContentTypeIdentifier{
			"htm":      htmlcontent.TypeIdentifier,
			"html":     htmlcontent.TypeIdentifier,
			"jsonld":   jsonldcontent.TypeIdentifier,
			"nq":       nquadscontent.TypeIdentifier,
			"nquads":   nquadscontent.TypeIdentifier,
			"nt":       ntriplescontent.TypeIdentifier,
			"ntriples": ntriplescontent.TypeIdentifier,
			"rdfjson":  rdfjsoncontent.TypeIdentifier,
			"rdfxml":   rdfxmlcontent.TypeIdentifier,
			"rj":       rdfjsoncontent.TypeIdentifier,
			"trig":     trigcontent.TypeIdentifier,
			"ttl":      turtlecontent.TypeIdentifier,
			"turtle":   turtlecontent.TypeIdentifier,
			"xhtml":    htmlcontent.TypeIdentifier,
			"xml":      rdfxmlcontent.TypeIdentifier,
			// dev
			"dev/null":    encodingtest.DiscardEncoderContentTypeIdentifier,
			"dev/quads":   encodingtest.QuadsEncoderContentTypeIdentifier,
			"dev/triples": encodingtest.TriplesEncoderContentTypeIdentifier,
		},
		MediaTypes: map[string]encoding.ContentTypeIdentifier{
			"application/ld+json":   jsonldcontent.TypeIdentifier,
			"application/n-quads":   nquadscontent.TypeIdentifier,
			"application/n-triples": ntriplescontent.TypeIdentifier,
			"application/rdf+json":  rdfjsoncontent.TypeIdentifier,
			"application/rdf+xml":   rdfxmlcontent.TypeIdentifier,
			"application/trig":      trigcontent.TypeIdentifier,
			"application/xhtml+xml": htmlcontent.TypeIdentifier,
			"text/html":             htmlcontent.TypeIdentifier,
			"text/turtle":           turtlecontent.TypeIdentifier,
			"text/xhtml+xml":        htmlcontent.TypeIdentifier,
		},
		FileExts: map[string]encoding.ContentTypeIdentifier{
			".htm":    htmlcontent.TypeIdentifier,
			".html":   htmlcontent.TypeIdentifier,
			".jsonld": jsonldcontent.TypeIdentifier,
			".nq":     nquadscontent.TypeIdentifier,
			".nt":     ntriplescontent.TypeIdentifier,
			".rdf":    rdfxmlcontent.TypeIdentifier,
			".trig":   trigcontent.TypeIdentifier,
			".ttl":    turtlecontent.TypeIdentifier,
		},
		MagicBytesResolvers: []encodingref.MagicBytesResolver{
			encodingref.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if rdfjsoncontent.MatchBytes(buf) {
					return rdfjsoncontent.TypeIdentifier, true
				}

				return "", false
			}),
			encodingref.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if jsonldcontent.MatchBytes(buf) {
					return jsonldcontent.TypeIdentifier, true
				}

				return "", false
			}),
			encodingref.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if htmlcontent.MatchBytes(buf) {
					return htmlcontent.TypeIdentifier, true
				}

				return "", false
			}),
			encodingref.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if rdfxmlcontent.MatchBytes(buf) {
					return rdfxmlcontent.TypeIdentifier, true
				}

				return "", false
			}),
			encodingref.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if htmlcontent.MatchBytesLax(buf) {
					return htmlcontent.TypeIdentifier, true
				}

				return "", false
			}),
		},
		Encodings: map[encoding.ContentTypeIdentifier]encodingref.RegistryEncoding{
			htmlcontent.TypeIdentifier: encodingHtml{
				jsonldDocumentLoader: opts.DocumentLoaderJSONLD,
			},
			jsonldcontent.TypeIdentifier: encodingJsonld{
				jsonldDocumentLoader: opts.DocumentLoaderJSONLD,
			},
			ntriplescontent.TypeIdentifier: encodingNtriples{},
			nquadscontent.TypeIdentifier:   encodingNquads{},
			rdfxmlcontent.TypeIdentifier:   encodingRdfxml{},
			rdfjsoncontent.TypeIdentifier:  encodingRdfjson{},
			trigcontent.TypeIdentifier:     encodingTrig{},
			turtlecontent.TypeIdentifier:   encodingTurtle{},
			// dev
			ctiDevHtmlInspector: devEncoding,
			encodingtest.DiscardEncoderContentTypeIdentifier: devEncoding,
			encodingtest.QuadsEncoderContentTypeIdentifier:   devEncoding,
			encodingtest.TriplesEncoderContentTypeIdentifier: devEncoding,
		},
	})
}
