package htmlcontent

import (
	"regexp"

	"github.com/dpb587/rdfkit-go/encoding"
)

const TypeIdentifier encoding.ContentTypeIdentifier = "public.html"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".html",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "xhtml+xml",
	},
}

var (
	reMatchHTML         = regexp.MustCompile(`^(<[^>]+>|\s)*<html[\s>]`)
	reMatchVocab        = regexp.MustCompile(`(<[\w]+\s+[^>]*vocab=")`)
	reMatchItemscope    = regexp.MustCompile(`(<[\w]+\s+[^>]*itemscope(\s|=""|>))`)
	reMatchJSONLDScript = regexp.MustCompile(`<script[^>]+type="application/ld\+json(\s*;[^"]+)?"`)
)

func MatchBytes(buf []byte) bool {
	if reMatchHTML.Match(buf) {
		return true
	}

	return false
}

func MatchBytesLax(buf []byte) bool {
	if reMatchVocab.Match(buf) {
		return true
	} else if reMatchItemscope.Match(buf) {
		return true
	} else if reMatchJSONLDScript.Match(buf) {
		return true
	}

	return false
}
