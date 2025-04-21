package cmdflags

import (
	"hash"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type EncodingInputHandle struct {
	Format   string
	ReadPath string
	Decoder  encoding.Decoder

	DecodedBase           []string
	DecodedPrefixMappings iriutil.PrefixMappingList

	readHasher hash.Hash
	reader     io.ReadCloser
}

func (h *EncodingInputHandle) HashSum() []byte {
	return h.readHasher.Sum(nil)
}

func (b *EncodingInputHandle) Close() error {
	if err := b.Decoder.Close(); err != nil {
		return err
	}

	if err := b.reader.Close(); err != nil {
		return err
	}

	return nil
}
