package nquads

import (
	"bytes"
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestWriter_write_IRIREF(t *testing.T) {
	for _, tc := range []struct {
		InputIRI         rdf.IRI
		InputOptionASCII bool
		OutputBytes      []byte
		OutputError      string
	}{
		{
			InputIRI:    `http://example.com/`,
			OutputBytes: []byte(`<http://example.com/>`),
		},
		{
			InputIRI:    `http://example.com/ `,
			OutputBytes: []byte(`<http://example.com/\u0020>`),
		},
		{
			InputIRI:    `http://example.com/>`,
			OutputBytes: []byte(`<http://example.com/\u003E>`),
		},
		{
			InputIRI:    `http://example.com/🐛`,
			OutputBytes: []byte(`<http://example.com/🐛>`),
		},
		{
			InputIRI:         `http://example.com/🐛`,
			InputOptionASCII: true,
			OutputBytes:      []byte(`<http://example.com/\U0001F41B>`),
		},
	} {
		t.Run(string(tc.InputIRI), func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := WriteIRI(buf, tc.InputIRI, tc.InputOptionASCII)
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
