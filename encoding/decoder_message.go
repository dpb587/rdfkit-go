package encoding

type DecoderMessage interface {
	GetDecoder() Decoder
}

//

type DecoderMessageWriter interface {
	WriteMessage(msg DecoderMessage)
}

//

type DecoderMessageWriterFunc func(msg DecoderMessage)

func (f DecoderMessageWriterFunc) WriteMessage(msg DecoderMessage) {
	f(msg)
}
