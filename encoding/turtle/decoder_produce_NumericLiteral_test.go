package turtle

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
)

func TestDecoder_CaptureNumericLiteral(t *testing.T) {
	for _, tc := range []struct {
		InputString     string
		OutputValue     string
		OutputValueType grammar.R
		Error           string
	}{
		// [spec, 2.5.2] // Numbers
		{
			InputString:     `-5`,
			OutputValue:     `-5`,
			OutputValueType: grammar.R_INTEGER,
		},
		{
			InputString:     `-5.0`,
			OutputValue:     `-5.0`,
			OutputValueType: grammar.R_DECIMAL,
		},
		{
			InputString:     `4.2E9`,
			OutputValue:     `4.2E9`,
			OutputValueType: grammar.R_DOUBLE,
		},
		//
		{
			InputString:     `+54`,
			OutputValue:     `+54`,
			OutputValueType: grammar.R_INTEGER,
		},
		{
			InputString:     `+54.04`,
			OutputValue:     `+54.04`,
			OutputValueType: grammar.R_DECIMAL,
		},
		{
			InputString:     `43.21e98`,
			OutputValue:     `43.21e98`,
			OutputValueType: grammar.R_DOUBLE,
		},
		{
			InputString:     `43.21E+98`,
			OutputValue:     `43.21E+98`,
			OutputValueType: grammar.R_DOUBLE,
		},
		{
			InputString:     `43.21E-98`,
			OutputValue:     `43.21E-98`,
			OutputValueType: grammar.R_DOUBLE,
		},
		{
			InputString:     `.5`,
			OutputValue:     `.5`,
			OutputValueType: grammar.R_DECIMAL,
		},
		// EOF
		{
			InputString: `-`,
			Error:       `token (NumericLiteral): offset 0x0: unexpected rune ('-')`,
		},
		{
			InputString: `+`,
			Error:       `token (NumericLiteral): offset 0x0: unexpected rune ('+')`,
		},
		{
			InputString:     `+5.`,
			OutputValue:     `+5`, // against grammar per turtle-syntax-number-08
			OutputValueType: grammar.R_INTEGER,
		},
		{
			InputString: `4E`,
			Error:       `token (DOUBLE): token (EXPONENT): offset 0x2: EOF`,
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

			token, err := s.produceNumericLiteral(r0)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.OutputValue, token.Decoded; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			} else if _e, _a := tc.OutputValueType, token.GrammarRule; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
