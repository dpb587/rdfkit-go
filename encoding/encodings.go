package encoding

import (
	"github.com/dpb587/rdfkit-go/rdfio"
)

type ContentMetadata struct {
	FileExt string

	MediaType  string
	Charset    string
	Parameters []string
}

type Decoder interface {
	rdfio.StatementIterator
}

type GraphDecoder interface {
	Decoder

	rdfio.GraphStatementIterator
}

type Encoder interface {
	rdfio.GraphWriter

	GetContentMetadata() ContentMetadata
}

type GraphEncoder interface {
	Encoder
}

type DatasetDecoder interface {
	Decoder

	rdfio.DatasetStatementIterator
}

type DatasetEncoder interface {
	rdfio.DatasetWriter

	GetContentMetadata() ContentMetadata
}
