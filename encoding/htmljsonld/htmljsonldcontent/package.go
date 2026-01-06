package htmljsonldcontent

import "github.com/dpb587/rdfkit-go/encoding"

const TypeIdentifier encoding.ContentTypeIdentifier = "public.html-json-ld"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".html",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "xhtml+xml",
	},
}
