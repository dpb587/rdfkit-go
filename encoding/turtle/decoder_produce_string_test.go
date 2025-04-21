package turtle

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
)

func TestDecoder_CaptureString(t *testing.T) {
	for _, tc := range []struct {
		InputString   string
		DecodedString string
		Error         string
	}{
		{
			InputString:   `""`,
			DecodedString: "",
		},
		{
			InputString:   `""""""`,
			DecodedString: "",
		},
		{
			InputString:   `''`,
			DecodedString: "",
		},
		{
			InputString:   `''''''`,
			DecodedString: "",
		},
		//
		{
			InputString:   `"hello"`,
			DecodedString: `hello`,
		},
		{
			InputString:   `"""hello"""`,
			DecodedString: `hello`,
		},
		{
			InputString:   `'hello'`,
			DecodedString: `hello`,
		},
		{
			InputString:   `'''hello'''`,
			DecodedString: `hello`,
		},
		//
		{
			InputString: `"hello\"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): offset 0x8: EOF`,
		},
		{
			InputString: `"""hello\"`,
			Error:       `token (String): token (STRING_LITERAL_LONG_QUOTE): offset 0xa: EOF`,
		},
		{
			InputString: `'hello\'`,
			Error:       `token (String): token (STRING_LITERAL_SINGLE_QUOTE): offset 0x8: EOF`,
		},
		{
			InputString: `'''hello\'`,
			Error:       `token (String): token (STRING_LITERAL_LONG_SINGLE_QUOTE): offset 0xa: EOF`,
		},
		//
		{
			InputString:   `"hello\\"`,
			DecodedString: `hello\`,
		},
		{
			InputString:   `"""hello\\"""`,
			DecodedString: `hello\`,
		},
		{
			InputString:   `'hello\\'`,
			DecodedString: `hello\`,
		},
		{
			InputString:   `'''hello\\'''`,
			DecodedString: `hello\`,
		},
		//
		{
			InputString:   `""X`,
			DecodedString: ``,
		},
		{
			InputString:   `''X`,
			DecodedString: ``,
		},
		//
		{
			InputString:   `"""hello"test""world"""`,
			DecodedString: `hello"test""world`,
		},
		{
			InputString:   `"""hello'test''world"""`,
			DecodedString: `hello'test''world`,
		},
		{
			InputString:   `'''hello'test''world'''`,
			DecodedString: `hello'test''world`,
		},
		{
			InputString:   `'''hello"test""world'''`,
			DecodedString: `hello"test""world`,
		},
		//
		{
			InputString: `"hello\X"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (ECHAR): offset 0x7: unexpected rune ('X')`,
		},
		{
			InputString: `"hello\u"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\u000"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0xb: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\uXXXX"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('X')`,
		},
		{
			InputString:   `"hello\u003E"`,
			DecodedString: `hello>`,
		},
		{
			InputString:   `"hello\u003e"`,
			DecodedString: `hello>`,
		},
		{
			InputString: `"hello\U"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\U0000000"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0xf: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\UXXXXXXXX"`,
			Error:       `token (String): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('X')`,
		},
		{
			InputString:   `"hello\U0001F41B"`,
			DecodedString: `hello🐛`,
		},
		{
			InputString:   `"hello\U0001f41b"`,
			DecodedString: `hello🐛`,
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

			token, err := s.produceString(r0)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.DecodedString, token.Decoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
