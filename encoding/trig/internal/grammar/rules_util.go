package grammar

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
)

func (r R) Err(err error) error {
	return encodingutil.WrapScanToken(r, err, nil)
}

func (r R) ErrCursorRange(err error, cr cursorio.OffsetRange) error {
	return encodingutil.WrapScanToken(r, err, cr)
}
