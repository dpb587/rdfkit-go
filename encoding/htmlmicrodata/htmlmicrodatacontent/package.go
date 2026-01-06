package htmlmicrodatacontent

import "github.com/dpb587/rdfkit-go/encoding"

const TypeIdentifier encoding.ContentTypeIdentifier = "public.html-microdata"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".html",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "xhtml+xml",
	},
}
