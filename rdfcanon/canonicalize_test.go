package rdfcanon_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cespare/permute/v2"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/quads"
	"github.com/dpb587/rdfkit-go/rdfcanon"
)

func TestCanonicalize_ContextCancellation(t *testing.T) {
	var statements rdf.QuadList // ~500k

	for range 4 {
		var bnSet = make([]rdf.BlankNode, 8)

		for j := range bnSet {
			bnSet[j] = rdf.NewBlankNode()
		}

		bnSetPermutations := permute.Slice(bnSet)

		for bnSetPermutations.Permute() {
			statements = append(statements,
				rdf.Quad{Triple: rdf.Triple{Subject: bnSet[0], Predicate: rdf.IRI("http://example.org/p"), Object: bnSet[1]}},
				rdf.Quad{Triple: rdf.Triple{Subject: bnSet[1], Predicate: rdf.IRI("http://example.org/p"), Object: bnSet[2]}},
				rdf.Quad{Triple: rdf.Triple{Subject: bnSet[2], Predicate: rdf.IRI("http://example.org/p"), Object: bnSet[3]}},
			)
		}
	}

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := rdfcanon.Canonicalize(ctx, quads.NewIterator(statements))
		if _a, _e := err, context.Canceled; !errors.Is(_a, _e) {
			t.Errorf("expected %v error, got: %v", _e, _a)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := rdfcanon.Canonicalize(ctx, quads.NewIterator(statements))
		if _a, _e := err, context.DeadlineExceeded; !errors.Is(_a, _e) {
			t.Errorf("expected %v error, got: %v", _e, _a)
		}
	})
}
