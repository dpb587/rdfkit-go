package ioutil

import "io"

type WriterCloser struct {
	io.Writer

	closers []func() error
}

func NewWriterCloser(reader io.Writer, closers ...func() error) *WriterCloser {
	return &WriterCloser{
		Writer:  reader,
		closers: closers,
	}
}

func (r *WriterCloser) Close() error {
	var errs []error

	for _, closer := range r.closers {
		err := closer()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
