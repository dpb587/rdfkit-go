package turtle

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/rdfdescription"
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
		{
			name: "list with single anon resource - single property",
			content: `<http://example.com/weight_Property>
	<http://www.w3.org/ns/shacl#name> "Weight"@en ;
	<http://www.w3.org/ns/shacl#or> (
		[ <http://www.w3.org/ns/shacl#class> <http://schema.org/Mass> ]
	) ;
	<http://www.w3.org/ns/shacl#path> <http://schema.org/weight> .
`,
		},
		{
			name: "list with single anon resource - multiple properties",
			content: `<http://example.com/weight_Property>
	<http://www.w3.org/ns/shacl#name> "Weight"@en ;
	<http://www.w3.org/ns/shacl#or> (
		[
			<http://example.com/name> "name"@en ;
			<http://www.w3.org/ns/shacl#class> <http://schema.org/Mass>
		]
	) ;
	<http://www.w3.org/ns/shacl#path> <http://schema.org/weight> .
`,
		},
		{
			name: "list with single simple value",
			content: `<http://example.com/test> <http://example.com/items> (
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

			resources := rdfdescription.NewResourceListBuilder()

			for decoder.Next() {
				resources.AddTriple(t.Context(), decoder.Triple())
			}

			if err := decoder.Err(); err != nil {
				t.Fatal(err)
			}

			err = resources.ToResourceWriter(ctx, encoder, rdfdescription.DefaultExportResourceOptions)
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
