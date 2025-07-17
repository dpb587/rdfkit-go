package ntriples

import (
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func TestDecoder_ScanIRIREF(t *testing.T) {
	for _, tc := range []struct {
		InputString string
		OutputIRI   rdf.IRI
		OutputError string
	}{
		{
			InputString: `<http://example.com/>`,
			OutputIRI:   "http://example.com/",
		},
		{
			InputString: `<http://example.com/\>`,
			OutputError: `token (IRIREF): offset 0x15: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\X>`,
			OutputError: `token (IRIREF): offset 0x15: unexpected rune ('X')`,
		},
		{
			InputString: `<http://example.com/\u>`,
			OutputError: `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\u000>`,
			OutputError: `token (IRIREF): token (UCHAR): offset 0x19: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\uXXXX>`,
			OutputError: `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('X')`,
		},
		{
			InputString: `<http://example.com/\u003E>`,
			OutputIRI:   "http://example.com/>",
		},
		{
			InputString: `<http://example.com/\u003e>`,
			OutputIRI:   "http://example.com/>",
		},
		{
			InputString: `<http://example.com/\U>`,
			OutputError: `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\U0000000>`,
			OutputError: `token (IRIREF): token (UCHAR): offset 0x1d: unexpected rune ('>')`,
		},
		{
			InputString: `<http://example.com/\UXXXXXXXX>`,
			OutputError: `token (IRIREF): token (UCHAR): offset 0x16: unexpected rune ('X')`,
		},
		{
			InputString: `<http://example.com/\U0001F41B>`,
			OutputIRI:   "http://example.com/🐛",
		},
		{
			InputString: `<http://example.com/\U0001f41b>`,
			OutputIRI:   "http://example.com/🐛",
		},
	} {
		t.Run(string(tc.InputString), func(t *testing.T) {
			s, err := NewDecoder(strings.NewReader(tc.InputString))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			sr, _ := s.buf.NextRune()

			outputIRI, _, err := s.captureOpenIRI(cursorio.DecodedRuneList{sr})
			if err == nil && len(tc.OutputError) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil && err.Error() != tc.OutputError {
				t.Errorf("unexpected error: %s", err)
			} else if !outputIRI.TermEquals(tc.OutputIRI) {
				t.Errorf("unexpected output: %s", outputIRI)
			}
		})
	}
}

func TestDecoder_ScanString(t *testing.T) {
	for _, tc := range []struct {
		InputString   string
		OutputLiteral rdf.Literal
		OutputError   string
	}{
		{
			InputString: `"hello"`,
			OutputLiteral: rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: "hello",
			},
		},
		{
			InputString: `"hello\"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): offset 0x8: EOF`,
		},
		{
			InputString: `"hello\\"`,
			OutputLiteral: rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: "hello\\",
			},
		},
		{
			InputString: `"hello\X"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): offset 0x7: unexpected rune ('X')`,
		},
		{
			InputString: `"hello\u"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\u000"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0xb: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\uXXXX"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('X')`,
		},
		{
			InputString: `"hello\u003E"`,
			OutputLiteral: rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: "hello>",
			},
		},
		{
			InputString: `"hello\u003e"`,
			OutputLiteral: rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: "hello>",
			},
		},
		{
			InputString: `"hello\U"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\U0000000"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0xf: unexpected rune ('"')`,
		},
		{
			InputString: `"hello\UXXXXXXXX"`,
			OutputError: `token (literal): token (STRING_LITERAL_QUOTE): token (UCHAR): offset 0x8: unexpected rune ('X')`,
		},
		{
			InputString: `"hello\U0001F41B"`,
			OutputLiteral: rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: "hello🐛",
			},
		},
		{
			InputString: `"hello\U0001f41b"`,
			OutputLiteral: rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: "hello🐛",
			},
		},
	} {
		t.Run(string(tc.InputString), func(t *testing.T) {
			s, err := NewDecoder(strings.NewReader(tc.InputString))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			sr, _ := s.buf.NextRune()

			outputLiteral, _, err := s.captureOpenLiteral(cursorio.DecodedRuneList{sr})
			if err == nil && len(tc.OutputError) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil && err.Error() != tc.OutputError {
				t.Errorf("unexpected error: %s", err)
			} else if !outputLiteral.TermEquals(tc.OutputLiteral) {
				t.Errorf("unexpected output: %#+v", outputLiteral)
			}
		})
	}
}
