package trig

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
)

func TestDecoder_CaptureIRIREF(t *testing.T) {
	for _, tc := range []struct {
		InputString        string
		OutputDecodedValue string
		Error              string
	}{
		{
			InputString:        `<http://example.com/>`,
			OutputDecodedValue: `http://example.com/`,
		},
		{
			InputString: `<http://example.com/\>`,
			Error:       `token (IRIREF): offset 0x15: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\X>`,
			Error:       `token (IRIREF): offset 0x15: unexpected rune ('X')`,
		},
		{
			InputString: `<http://example.com/\u>`,
			Error:       `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\u000>`,
			Error:       `token (IRIREF): token (UCHAR): offset 0x19: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\uXXXX>`,
			Error:       `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('X')`,
		},
		{
			InputString:        `<http://example.com/\u003E>`,
			OutputDecodedValue: `http://example.com/>`,
		},
		{
			InputString:        `<http://example.com/\u003e>`,
			OutputDecodedValue: `http://example.com/>`,
		},
		{
			InputString: `<http://example.com/\U>`,
			Error:       `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\U0000000>`,
			Error:       `token (IRIREF): token (UCHAR): offset 0x1d: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\U10000000>`,
			Error:       `token (IRIREF): token (UCHAR): exceeds maximum unicode code point`,
		},
		{
			InputString: `<http://example.com/\UXXXXXXXX>`,
			Error:       `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('X')`,
		},
		{
			InputString:        `<http://example.com/\U0001F41B>`,
			OutputDecodedValue: `http://example.com/🐛`,
		},
		{
			InputString:        `<http://example.com/\U0001f41b>`,
			OutputDecodedValue: `http://example.com/🐛`,
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

			token, err := s.produceIRIREF(r0)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.OutputDecodedValue, token.Decoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
