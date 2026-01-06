package nquadscontent

import "github.com/dpb587/rdfkit-go/encoding"

const TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.n-quads"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".nq",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "n-quads",
	},
}
