package trigcontent

import "github.com/dpb587/rdfkit-go/encoding"

const TypeIdentifier encoding.ContentTypeIdentifier = "org.w3.trig"

var DefaultMetadata = encoding.ContentMetadata{
	FileExt: ".trig",
	MediaType: encoding.ContentMediaType{
		Type:    "application",
		Subtype: "trig",
	},
}
