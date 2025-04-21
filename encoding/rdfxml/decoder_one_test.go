package rdfxml

import (
	"strings"
	"testing"
)

func TestOne(t *testing.T) {
	r, err := NewDecoder(
		strings.NewReader(`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:schema="https://schema.org/">
  <schema:WebSite>
    <schema:name>Named Graph</schema:name>
    <schema:url rdf:resource="https://www.namedgraph.com/" />
  </schema:WebSite>
</rdf:RDF>
`),
		DecoderConfig{}.
			SetCaptureTextOffsets(true),
	)
	if err != nil {
		t.Fatalf("failed to create decoder: %v", err)
	}

	for r.Next() {
		// true
	}
	if r.Err() != nil {
		panic(r.Err())
	}

}
