package jsonldcontent

import (
	"regexp"

	"github.com/dpb587/rdfkit-go/encoding"
)

const TypeIdentifier encoding.ContentTypeIdentifier = "org.json-ld.document"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".jsonld",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "ld+json",
	},
}

var (
	reMatchArrayObject = regexp.MustCompile(`^\s*\[\s*\{`)
	reMatchObject      = regexp.MustCompile(`^\s*\{`)
)

func MatchBytes(buf []byte) bool {
	if reMatchArrayObject.Match(buf) {
		return true
	} else if reMatchObject.Match(buf) {
		return true
	}

	return false
}
