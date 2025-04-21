package turtle

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
)

func TestDecoder_CaptureLANGTAG(t *testing.T) {
	for _, tc := range []struct {
		InputString     string
		OutputLangValue string
		Error           string
	}{
		{
			InputString:     `@en`,
			OutputLangValue: `en`,
		},
		{
			InputString:     `@en^^`,
			OutputLangValue: `en`,
		},
		{
			InputString:     `@fr-be`,
			OutputLangValue: `fr-be`,
		},
		{
			InputString: `@`,
			Error:       "token (LANGTAG): offset 0x1: EOF",
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

			token, err := s.produceLANGTAG(r0)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.OutputLangValue, token.Decoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
