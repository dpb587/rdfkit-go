package html

type DocumentInfo struct {
	// Location is the original location provided when parsing the document.
	//
	// BaseURL should typically be used when resolving relative URLs of content.
	Location string

	// BaseURL is the URL from the first <base href> tag, or the original location.
	BaseURL string

	// HasNodeMetadata indicates whether node-level metadata with text offsets is available.
	HasNodeMetadata bool
}
