package cmdflags

import (
	"hash"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
)

type EncodingOutputHandle struct {
	Format    string
	WritePath string
	Encoder   encoding.Encoder

	writeHasher hash.Hash
	writer      io.WriteCloser
}

func (h *EncodingOutputHandle) HashSum() []byte {
	return h.writeHasher.Sum(nil)
}

func (b *EncodingOutputHandle) Close() error {
	if err := b.Encoder.Close(); err != nil {
		return err
	}

	if err := b.writer.Close(); err != nil {
		return err
	}

	return nil
}
