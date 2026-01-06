package rdfio

import (
	"net/http"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/html/htmlcontent"
	"github.com/dpb587/rdfkit-go/encoding/html/htmldefaults/htmldefaultsrdfio"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldcontent"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldrdfio"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadsrdfio"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/ntriplescontent"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/ntriplesrdfio"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson/rdfjsoncontent"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson/rdfjsonrdfio"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml/rdfxmlcontent"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml/rdfxmlrdfio"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigcontent"
	"github.com/dpb587/rdfkit-go/encoding/trig/trigrdfio"
	"github.com/dpb587/rdfkit-go/encoding/turtle/turtlecontent"
	"github.com/dpb587/rdfkit-go/encoding/turtle/turtlerdfio"
	"github.com/dpb587/rdfkit-go/rdfio/fileresource"
	"github.com/dpb587/rdfkit-go/rdfio/httpresource"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type RegistryOptions struct {
	HttpClient           *http.Client
	DocumentLoaderJSONLD jsonldtype.DocumentLoader
}

func NewRegistry(opts RegistryOptions) rdfiotypes.Registry {
	return rdfiotypes.Registry{
		Aliases: map[string]encoding.ContentTypeIdentifier{
			"dev/null":    encodingtest.DiscardEncoderContentTypeIdentifier,
			"dev/quads":   encodingtest.QuadsEncoderContentTypeIdentifier,
			"dev/triples": encodingtest.TriplesEncoderContentTypeIdentifier,
			"htm":         htmlcontent.TypeIdentifier,
			"html":        htmlcontent.TypeIdentifier,
			"jsonld":      jsonldcontent.TypeIdentifier,
			"n-quads":     nquadscontent.TypeIdentifier,
			"n-triples":   ntriplescontent.TypeIdentifier,
			"nq":          nquadscontent.TypeIdentifier,
			"nquads":      nquadscontent.TypeIdentifier,
			"nt":          ntriplescontent.TypeIdentifier,
			"ntriples":    ntriplescontent.TypeIdentifier,
			"rdf-json":    rdfjsoncontent.TypeIdentifier,
			"rdf-xml":     rdfxmlcontent.TypeIdentifier,
			"rdfjson":     rdfjsoncontent.TypeIdentifier,
			"rdfxml":      rdfxmlcontent.TypeIdentifier,
			"rj":          rdfjsoncontent.TypeIdentifier,
			"trig":        trigcontent.TypeIdentifier,
			"ttl":         turtlecontent.TypeIdentifier,
			"turtle":      turtlecontent.TypeIdentifier,
			"xhtml":       htmlcontent.TypeIdentifier,
			"xml":         rdfxmlcontent.TypeIdentifier,
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
			".rj":     rdfjsoncontent.TypeIdentifier,
			".trig":   trigcontent.TypeIdentifier,
			".ttl":    turtlecontent.TypeIdentifier,
			".xhtml":  htmlcontent.TypeIdentifier,
		},
		MagicBytesResolvers: []rdfiotypes.MagicBytesResolver{
			rdfiotypes.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if rdfjsoncontent.MatchBytes(buf) {
					return rdfjsoncontent.TypeIdentifier, true
				}

				return "", false
			}),
			rdfiotypes.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if jsonldcontent.MatchBytes(buf) {
					return jsonldcontent.TypeIdentifier, true
				}

				return "", false
			}),
			rdfiotypes.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if htmlcontent.MatchBytes(buf) {
					return htmlcontent.TypeIdentifier, true
				}

				return "", false
			}),
			rdfiotypes.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if rdfxmlcontent.MatchBytes(buf) {
					return rdfxmlcontent.TypeIdentifier, true
				}

				return "", false
			}),
			rdfiotypes.MagicBytesResolverFunc(func(buf []byte) (encoding.ContentTypeIdentifier, bool) {
				if htmlcontent.MatchBytesLax(buf) {
					return htmlcontent.TypeIdentifier, true
				}

				return "", false
			}),
		},
		ResourceManagers: []rdfiotypes.ResourceManager{
			httpresource.NewManager(http.DefaultClient),
			fileresource.NewManager(),
		},
		DecoderManagers: map[encoding.ContentTypeIdentifier]rdfiotypes.DecoderManager{
			htmlcontent.TypeIdentifier:     htmldefaultsrdfio.NewDecoder(opts.DocumentLoaderJSONLD),
			jsonldcontent.TypeIdentifier:   jsonldrdfio.NewDecoder(opts.DocumentLoaderJSONLD),
			ntriplescontent.TypeIdentifier: ntriplesrdfio.NewDecoder(),
			nquadscontent.TypeIdentifier:   nquadsrdfio.NewDecoder(),
			rdfxmlcontent.TypeIdentifier:   rdfxmlrdfio.NewDecoder(),
			rdfjsoncontent.TypeIdentifier:  rdfjsonrdfio.NewDecoder(),
			trigcontent.TypeIdentifier:     trigrdfio.NewDecoder(),
			turtlecontent.TypeIdentifier:   turtlerdfio.NewDecoder(),
		},
		EncoderManagers: map[encoding.ContentTypeIdentifier]rdfiotypes.EncoderManager{
			ntriplescontent.TypeIdentifier:                   ntriplesrdfio.NewEncoder(),
			nquadscontent.TypeIdentifier:                     nquadsrdfio.NewEncoder(),
			rdfjsoncontent.TypeIdentifier:                    rdfjsonrdfio.NewEncoder(),
			turtlecontent.TypeIdentifier:                     turtlerdfio.NewEncoder(),
			ctiDevHtmlInspector:                              encodingDevHtmlInspector{},
			encodingtest.DiscardEncoderContentTypeIdentifier: encodingDevDiscard{},
			encodingtest.QuadsEncoderContentTypeIdentifier:   encodingDevQuads{},
			encodingtest.TriplesEncoderContentTypeIdentifier: encodingDevTriples{},
		},
	}
}
