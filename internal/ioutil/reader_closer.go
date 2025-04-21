package ioutil

import "io"

type ReaderCloser struct {
	io.Reader

	closers []func() error
}

func NewReaderCloser(reader io.Reader, closers ...func() error) *ReaderCloser {
	return &ReaderCloser{
		Reader:  reader,
		closers: closers,
	}
}

func (r *ReaderCloser) Close() error {
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
