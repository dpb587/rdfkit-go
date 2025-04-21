package nquads

import (
	"bytes"
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func TestNewStringLiteralQuote(t *testing.T) {
	for _, tc := range []struct {
		InputString      string
		InputOptionASCII bool
		OutputBytes      []byte
		OutputError      string
	}{
		{
			InputString: `hello`,
			OutputBytes: []byte(`"hello"`),
		},
		{
			InputString: `hello `,
			OutputBytes: []byte(`"hello "`),
		},
		{
			InputString: `hello"`,
			OutputBytes: []byte(`"hello\""`),
		},
		{
			InputString: `hello🐛`,
			OutputBytes: []byte(`"hello🐛"`),
		},
		{
			InputString:      `hello🐛`,
			InputOptionASCII: true,
			OutputBytes:      []byte(`"hello\U0001F41B"`),
		},
	} {
		t.Run(string(tc.InputString), func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := WriteLiteral(buf, rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: tc.InputString,
			}, tc.InputOptionASCII)
			if err == nil && len(tc.OutputError) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil && err.Error() != tc.OutputError {
				t.Errorf("unexpected error: %s", err)
			} else if _e, _a := tc.OutputBytes, buf.Bytes(); !bytes.Equal(_e, _a) {
				t.Errorf("unexpected output: %s", _a)
			}
		})
	}
}
