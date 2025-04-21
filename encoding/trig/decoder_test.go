package trig

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDecoder(t *testing.T) {
	r, err := NewDecoder(
		strings.NewReader(`@base <http://example.org/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix rel: <http://www.perceive.net/schemas/relationship/> .

<#green-goblin>
    rel:enemyOf <#spiderman> ;
    a foaf:Person ;    # in the context of the Marvel universe
    foaf:name "Green Goblin" .

<#spiderman>
    rel:enemyOf <#green-goblin> ;
    a foaf:Person ;
    foaf:name "Spiderman", "Человек-паук"@ru .`),
		DecoderConfig{}.
			SetCaptureTextOffsets(true),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for r.Next() {
		fmt.Fprintf(os.Stderr, "%v\n", r.GetTriple())
	}

	if err := r.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
