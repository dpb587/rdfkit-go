package ntriplescontent

import "github.com/dpb587/rdfkit-go/encoding"

const TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.n-triples"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".nt",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "n-triples",
	},
}
