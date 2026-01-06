package encoding

type ContentTypeIdentifier string

type ContentMetadata struct {
	FileExt   string
	MediaType ContentMediaType
}

type ContentMediaType struct {
	Type       string
	Subtype    string
	Parameters map[string]string
}
