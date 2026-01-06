package rdfxmlcontent

import (
	"regexp"

	"github.com/dpb587/rdfkit-go/encoding"
)

const TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.rdf-xml"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".rdf",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "rdf+xml",
	},
}

var (
	reMatchDeclaration = regexp.MustCompile(`^(<[^>]+>|\s)*<\?xml`)
	reMatchDefaultRDF  = regexp.MustCompile(`<rdf:RDF `)
)

func MatchBytes(buf []byte) bool {
	if reMatchDeclaration.Match(buf) {
		return true
	} else if reMatchDefaultRDF.Match(buf) {
		return true
	}

	return false
}
