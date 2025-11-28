package turtle

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionio"
)

func TestEncoder_Basic(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "rdf:type",
			content: `<http://example.com/test> a <http://example.com/Thing> .
`,
		},
		{
			name: "rdf:type multiple",
			content: `<http://example.com/test> a
	<http://example.com/Thing> ,
	<http://example.com/Thing2> .
`,
		},
		{
			name: "object anon",
			content: `<http://example.com/test> <http://example.com/predicate> [
	a <http://example.com/Thing2> ;
	<http://example.com/hasValue> "value"
] .
`,
		},
		{
			name: "sibling object anon",
			content: `<http://example.com/test> <http://example.com/predicate>
	<http://example.com/Thing> ,
	[
		a <http://example.com/Thing2> ;
		<http://example.com/hasValue> "value"
	] .
`,
		},
		{
			name: "simple list",
			content: `<http://example.com/test>
	a <http://example.com/Thing> ;
	<http://example.com/hasItems> (
		"item1"
		"item2"
	) .
`,
		},
		{
			name: "empty list",
			content: `<http://example.com/test> <http://example.com/hasItems> <http://www.w3.org/1999/02/22-rdf-syntax-ns#nil> .
`,
		},
		{
			name: "single item list",
			content: `<http://example.com/test> <http://example.com/hasItems> (
	"item1"
) .
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			var buf bytes.Buffer

			decoder, err := NewDecoder(strings.NewReader(tt.content))
			if err != nil {
				t.Fatal(err)
			}

			encoder, err := NewEncoder(&buf, EncoderConfig{})
			if err != nil {
				t.Fatal(err)
			}

			resources := rdfdescriptionio.NewBufferedGraphEncoder(encoder)

			for decoder.Next() {
				err := resources.PutTriple(ctx, decoder.GetTriple())
				if err != nil {
					t.Fatal(err)
				}
			}

			if err := decoder.Err(); err != nil {
				t.Fatal(err)
			}

			err = resources.Close()
			if err != nil {
				t.Fatal(err)
			}

			err = encoder.Close()
			if err != nil {
				t.Fatal(err)
			}

			got := buf.String()
			if got != tt.content {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.content, got)
			}
		})
	}
}
