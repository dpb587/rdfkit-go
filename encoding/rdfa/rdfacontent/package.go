package rdfacontent

import "github.com/dpb587/rdfkit-go/encoding"

var TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.rdfa"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".html",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "xhtml+xml",
	},
}
