package encoding

import "github.com/dpb587/rdfkit-go/rdf"

type Decoder interface {
	rdf.StatementIterator

	GetContentTypeIdentifier() ContentTypeIdentifier
}

type Encoder interface {
	Close() error

	GetContentMetadata() ContentMetadata
	GetContentTypeIdentifier() ContentTypeIdentifier
}
