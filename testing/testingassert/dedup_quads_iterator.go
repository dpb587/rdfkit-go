package testingassert

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/dpb587/rdfkit-go/rdf"
)

type dedupQuadsIterator struct {
	iter rdf.QuadIterator
	seen map[string]struct{}
}

var _ rdf.QuadIterator = &dedupQuadsIterator{}

func newDedupQuadsIterator(iter rdf.QuadIterator) *dedupQuadsIterator {
	return &dedupQuadsIterator{
		iter: iter,
		seen: make(map[string]struct{}),
	}
}

func (d *dedupQuadsIterator) Next() bool {
	for {
		if !d.iter.Next() {
			return false
		}

		key := d.hash(d.iter.Quad())

		if _, exists := d.seen[key]; exists {
			continue
		}

		d.seen[key] = struct{}{}

		return true
	}
}

func (d *dedupQuadsIterator) Statement() rdf.Statement {
	return d.iter.Statement()
}

func (d *dedupQuadsIterator) Quad() rdf.Quad {
	return d.iter.Quad()
}

func (d *dedupQuadsIterator) Err() error {
	return d.iter.Err()
}

func (d *dedupQuadsIterator) Close() error {
	return d.iter.Close()
}

func (d *dedupQuadsIterator) hash(q rdf.Quad) string {
	h := sha256.New()
	d.hashTerm(h, q.Triple.Subject)
	d.hashTerm(h, q.Triple.Predicate)
	d.hashTerm(h, q.Triple.Object)
	d.hashTerm(h, q.GraphName)

	return hex.EncodeToString(h.Sum(nil)[0:16])
}

func (d *dedupQuadsIterator) hashTerm(h hash.Hash, t rdf.Term) {
	if t == nil {
		h.Write([]byte("nil\n"))

		return
	} else if bn, ok := t.(rdf.BlankNode); ok {
		fmt.Fprintf(h, "rdf.NewBlankNodeWithIdentifier(%#+v)\n", bn.GetBlankNodeIdentifier())

		return
	}

	fmt.Fprintf(h, "%#+v\n", t)
}
