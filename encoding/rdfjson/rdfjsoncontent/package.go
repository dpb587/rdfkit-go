package rdfjsoncontent

import (
	"regexp"

	"github.com/dpb587/rdfkit-go/encoding"
)

const TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.rdf-json"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".rj",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "rdf+json",
	},
}

var reMatchSyntax = regexp.MustCompile(`^\s*\{\s*"[^"]+"\s*:\s*\{\s*"[^"]+"\s*:\s*\[\s*\{\s*"(datatype|lang|type|value)"`)

func MatchBytes(buf []byte) bool {
	if reMatchSyntax.Match(buf) {
		return true
	}

	return false
}
