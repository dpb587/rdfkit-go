package turtlecontent

import "github.com/dpb587/rdfkit-go/encoding"

const TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.turtle"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".ttl",
	MediaType: encoding.ContentMediaType{
		Type:    "text",
		Subtype: "turtle",
		Parameters: map[string]string{
			"charset": "utf-8",
		},
	},
}
