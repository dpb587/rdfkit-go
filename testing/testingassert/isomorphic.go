package testingassert

import (
	"bytes"
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/quads"
	"github.com/dpb587/rdfkit-go/rdfcanon"
)

type testingT interface {
	Context() context.Context
	Fatalf(format string, args ...any)
}

func IsomorphicGraphs(t testingT, expected, actual rdf.TripleList) {
	IsomorphicDatasets(t, expected.AsQuads(nil), actual.AsQuads(nil))
}

func IsomorphicDatasets(t testingT, expected, actual rdf.QuadList) {
	expectedCanonical, err := rdfcanon.Canonicalize(t.Context(), newDedupQuadsIterator(quads.NewIterator(expected)))
	if err != nil {
		t.Fatalf("canonicalize expected: %v", err)
	}

	actualCanonical, err := rdfcanon.Canonicalize(t.Context(), newDedupQuadsIterator(quads.NewIterator(actual)))
	if err != nil {
		t.Fatalf("canonicalize actual: %v", err)
	}

	var expectedBuf bytes.Buffer

	_, err = expectedCanonical.WriteTo(&expectedBuf)
	if err != nil {
		t.Fatalf("write expected: %v", err)
	}

	var actualBuf bytes.Buffer

	_, err = actualCanonical.WriteTo(&actualBuf)
	if err != nil {
		t.Fatalf("write actual: %v", err)
	}

	if !bytes.Equal(expectedBuf.Bytes(), actualBuf.Bytes()) {
		t.Fatalf(
			"expected does not match actual\n\n=== EXPECTED\n%s\n\n=== ACTUAL\n%s",
			expectedBuf.String(),
			actualBuf.String(),
		)
	}
}
