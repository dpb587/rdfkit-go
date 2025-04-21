package rdfa

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/html"
)

func TestOne(t *testing.T) {
	d, err := html.ParseDocument(strings.NewReader(`<div vocab="https://schema.org/" typeof="Person">
  <div property="potentialAction" typeof="FindAction">
    <span property="query-input">name=hello placeholder=none</span>
</div>
</div>`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	r, err := NewDecoder(d)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for r.Next() {
		s := r.GetStatement()
		fmt.Fprintf(os.Stderr, "%#+v\n", s)
	}

	if r.Err() != nil {
		t.Fatalf("unexpected error: %s", r.Err())
	}
}
