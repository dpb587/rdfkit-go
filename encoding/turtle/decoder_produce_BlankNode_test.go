package turtle

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
)

func TestDecoder_CaptureBlankNode(t *testing.T) {
	for _, tc := range []struct {
		InputString           string
		OutputIdentifierValue string
		Error                 string
	}{
		{
			InputString:           `_:b`,
			OutputIdentifierValue: `b`,
		},
		{
			InputString:           `_:b0`,
			OutputIdentifierValue: `b0`,
		},
		{
			InputString: `_:`,
			Error:       `token (BLANK_NODE_LABEL): offset 0x2: EOF`,
		},
		{
			InputString: `_:\n`,
			Error:       `token (BLANK_NODE_LABEL): offset 0x2: unexpected rune ('\\')`,
		},
	} {
		t.Run(string(tc.InputString), func(t *testing.T) {
			s := &Decoder{
				buf: cursorioutil.NewRuneBuffer(strings.NewReader(tc.InputString)),
			}

			r0, err := s.buf.NextRune()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			token, err := s.produceBlankNode(r0)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.OutputIdentifierValue, token.Decoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
