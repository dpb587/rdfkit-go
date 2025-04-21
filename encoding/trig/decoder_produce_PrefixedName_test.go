package trig

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
)

func TestDecoder_CapturePrefixedName(t *testing.T) {
	for _, tc := range []struct {
		InputString             string
		OutputNamespaceValue    string
		OutputLocalDecodedValue string
		Error                   string
	}{
		// [spec, 2.4] // IRIs
		{
			InputString:          `somePrefix:`,
			OutputNamespaceValue: `somePrefix`,
		},
		{
			InputString:             `leg:3032571`,
			OutputNamespaceValue:    `leg`,
			OutputLocalDecodedValue: `3032571`,
		},
		{
			InputString:             `isbn13:9780136019701`,
			OutputNamespaceValue:    `isbn13`,
			OutputLocalDecodedValue: `9780136019701`,
		},
		{
			InputString:             `og:video:height`,
			OutputNamespaceValue:    `og`,
			OutputLocalDecodedValue: `video:height`,
		},
		{
			InputString:             `wgs:lat\-long`,
			OutputNamespaceValue:    `wgs`,
			OutputLocalDecodedValue: `lat-long`,
		},
		{
			InputString:             `:`,
			OutputNamespaceValue:    ``,
			OutputLocalDecodedValue: ``,
		},
		//
		{
			InputString:             `somePrefix: `,
			OutputNamespaceValue:    `somePrefix`,
			OutputLocalDecodedValue: ``,
		},
		{
			InputString:             `leg:3032571 `,
			OutputNamespaceValue:    `leg`,
			OutputLocalDecodedValue: `3032571`,
		},
		{
			InputString:             `isbn13:9780136019701 `,
			OutputNamespaceValue:    `isbn13`,
			OutputLocalDecodedValue: `9780136019701`,
		},
		{
			InputString:             `og:video:height `,
			OutputNamespaceValue:    `og`,
			OutputLocalDecodedValue: `video:height`,
		},
		{
			InputString:             `wgs:lat\-long `,
			OutputNamespaceValue:    `wgs`,
			OutputLocalDecodedValue: `lat-long`,
		},
		{
			InputString:             `: `,
			OutputNamespaceValue:    ``,
			OutputLocalDecodedValue: ``,
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

			token, err := s.producePrefixedName(r0)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.OutputNamespaceValue, token.NamespaceDecoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			} else if _e, _a := tc.OutputLocalDecodedValue, token.LocalDecoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
