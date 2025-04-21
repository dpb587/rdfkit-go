package encoding

import "io"

type GraphFactory interface {
	NewGraphDecoder(r io.Reader) (GraphDecoder, error)
	NewGraphEncoder(w io.Writer) (GraphEncoder, error)

	GetGraphEncoderContentMetadata() ContentMetadata
}

type DatasetFactory interface {
	GraphFactory

	NewDatasetDecoder(r io.Reader) (DatasetDecoder, error)
	NewDatasetEncoder(w io.Writer) (DatasetEncoder, error)

	GetDatasetEncoderContentMetadata() ContentMetadata
}
