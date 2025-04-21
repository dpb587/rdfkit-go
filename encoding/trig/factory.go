package trig

// import (
// 	"io"

// 	"github.com/dpb587/rdfkit-go/encoding"
// )

// type Factory struct {
// 	writerOpts EncoderOptions
// 	readerOpts DecoderOptions
// }

// var _ encoding.GraphFactory = &Factory{}

// func NewFactory(writerOpts EncoderOptions, readerOpts DecoderOptions) *Factory {
// 	return &Factory{
// 		writerOpts: writerOpts,
// 		readerOpts: readerOpts,
// 	}
// }

// func (e *Factory) NewGraphEncoder(w io.Writer) (encoding.GraphEncoder, error) {
// 	return NewEncoder(w, e.writerOpts), nil
// }

// func (e *Factory) NewGraphDecoder(r io.Reader) (encoding.GraphDecoder, error) {
// 	return NewDecoder(r, e.readerOpts), nil
// }

// func (e *Factory) GetGraphEncoderContentMetadata() encoding.ContentMetadata {
// 	return encoding.ContentMetadata{
// 		FileExt:   ".trig",
// 		MediaType: "application/trig",
// 		Charset:   "utf-8",
// 	}
// }
