package encodingdefaults

import (
	"io"

	"github.com/dpb587/rdfkit-go/x/encodingref"
)

func wrapReader(r io.Reader, opts encodingref.DecoderOptions) io.Reader {
	if opts.Tee != nil {
		r = io.TeeReader(r, opts.Tee)
	}

	if opts.Hasher != nil {
		r = io.TeeReader(r, opts.Hasher)
	}

	return r
}

func wrapWriter(w io.Writer, opts encodingref.EncoderOptions) io.Writer {
	if opts.Tee != nil {
		w = io.MultiWriter(w, opts.Tee)
	}

	if opts.Hasher != nil {
		w = io.MultiWriter(w, opts.Hasher)
	}

	return w
}
