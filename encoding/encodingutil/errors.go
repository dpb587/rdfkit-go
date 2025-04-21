package encodingutil

import (
	"fmt"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
)

type SyntaxDefinition interface {
	String() string
}

//

type TokenWrap struct {
	Token       SyntaxDefinition
	OffsetRange cursorio.OffsetRange
	Err         error
}

func WrapScanToken(t SyntaxDefinition, err error, cr cursorio.OffsetRange) TokenWrap {
	return TokenWrap{
		Token:       t,
		OffsetRange: cr,
		Err:         err,
	}
}

func (e TokenWrap) Error() string {
	s := &strings.Builder{}
	s.WriteString("token (" + e.Token.String())

	if e.OffsetRange != nil {
		s.WriteString("; offset=" + e.OffsetRange.OffsetRangeString())
	}

	s.WriteString("): ")
	fmt.Fprintf(s, "%v", e.Err)

	return s.String()
}

func (err TokenWrap) As(target interface{}) bool {
	if err.OffsetRange != nil {
		if tt, ok := target.(*cursorio.OffsetRangeError); ok {
			tt.OffsetRange = err.OffsetRange
			tt.Err = err

			return true
		}
	}

	return false
}

func (e TokenWrap) Unwrap() error {
	return e.Err
}
