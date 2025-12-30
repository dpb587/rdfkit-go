package testingutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/rdf/quads"
	"github.com/dpb587/rdfkit-go/rdfcanon"
)

type DebugRdfioBuilder struct {
	bufferByName map[string][]byte
}

func NewDebugRdfioBuilderFromEnv(t *testing.T) *DebugRdfioBuilder {
	bundle := &DebugRdfioBuilder{
		bufferByName: make(map[string][]byte),
	}

	if fhPath := os.Getenv("TESTING_DEBUG_RDFIO_OUTPUT"); len(fhPath) > 0 {
		t.Cleanup(func() {
			fh, err := os.OpenFile(fhPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				t.Fatalf("open debug file: %v", err)
			}

			defer fh.Close()

			_, err = bundle.WriteTo(fh)
			if err != nil {
				t.Errorf("close debug rdfio writer: %v", err)
			}
		})
	}

	return bundle
}

func (be *DebugRdfioBuilder) WriteTo(w io.Writer) (int64, error) {
	keys := slices.Collect(maps.Keys(be.bufferByName))
	slices.SortFunc(keys, strings.Compare)

	for nameIdx, name := range keys {
		if nameIdx > 0 {
			w.Write([]byte("\n"))
		}

		fmt.Fprintf(w, "=== %s\n", name)
		fmt.Fprintf(w, "\n")

		w.Write(be.bufferByName[name])
	}

	return 0, nil
}

type quadStatementWrapper struct {
	Statement      encodingtest.QuadStatement
	CanonizedIndex int
}

func (qsw quadStatementWrapper) EarliestTextOffset() int64 {
	var b cursorio.ByteOffset = math.MaxInt64

	for _, v := range qsw.Statement.TextOffsets {
		if v.From.Byte < b {
			b = v.From.Byte
		}
	}

	return int64(b)
}

func (be *DebugRdfioBuilder) PutQuadsBundle(name string, statements encodingtest.QuadStatementList) error {
	ctx := context.Background()

	// TODO statements.NewQuadIterator()
	canonicalized, err := rdfcanon.Canonicalize(quads.NewIterator(statements.AsQuads()))
	if err != nil {
		return fmt.Errorf("canonicalize: %v", err)
	}

	canonicalizedIter := canonicalized.NewIterator()

	var wrappedStatements []quadStatementWrapper

	for canonicalizedIter.Next() {
		wrappedStatements = append(wrappedStatements, quadStatementWrapper{
			Statement:      statements[canonicalizedIter.OriginalQuadIndex()],
			CanonizedIndex: len(wrappedStatements),
		})
	}

	slices.SortFunc(wrappedStatements, func(a, b quadStatementWrapper) int {
		ao, bo := a.EarliestTextOffset(), b.EarliestTextOffset()

		if ao < bo {
			return -1
		} else if ao > bo {
			return 1
		}

		return b.CanonizedIndex - a.CanonizedIndex
	})

	buf := &bytes.Buffer{}

	e := encodingtest.NewQuadsEncoder(buf, encodingtest.QuadsEncoderOptions{})

	for _, ws := range wrappedStatements {
		err := e.AddQuadStatement(ctx, ws.Statement)
		if err != nil {
			return err
		}
	}

	e.Close()

	be.bufferByName[name] = buf.Bytes()

	return nil
}

type tripleStatementWrapper struct {
	Statement      encodingtest.TripleStatement
	CanonizedIndex int
}

func (qsw tripleStatementWrapper) EarliestTextOffset() int64 {
	var b cursorio.ByteOffset = math.MaxInt64

	for _, v := range qsw.Statement.TextOffsets {
		if v.From.Byte < b {
			b = v.From.Byte
		}
	}

	return int64(b)
}

func (be *DebugRdfioBuilder) PutTriplesBundle(name string, statements encodingtest.TripleStatementList) error {
	ctx := context.Background()

	// TODO statements.NewTripleIterator()
	canonicalized, err := rdfcanon.Canonicalize(quads.NewIterator(statements.AsTriples().AsQuads(nil)))
	if err != nil {
		return fmt.Errorf("canonicalize: %v", err)
	}

	canonicalizedIter := canonicalized.NewIterator()

	var wrappedStatements []tripleStatementWrapper

	for canonicalizedIter.Next() {
		wrappedStatements = append(wrappedStatements, tripleStatementWrapper{
			Statement:      statements[canonicalizedIter.OriginalQuadIndex()],
			CanonizedIndex: len(wrappedStatements),
		})
	}

	slices.SortFunc(wrappedStatements, func(a, b tripleStatementWrapper) int {
		ao, bo := a.EarliestTextOffset(), b.EarliestTextOffset()

		if ao < bo {
			return -1
		} else if ao > bo {
			return 1
		}

		return b.CanonizedIndex - a.CanonizedIndex
	})

	buf := &bytes.Buffer{}

	e := encodingtest.NewTriplesEncoder(buf, encodingtest.TriplesEncoderOptions{})

	for _, ws := range wrappedStatements {
		err := e.AddTripleStatement(ctx, ws.Statement)
		if err != nil {
			return err
		}
	}

	e.Close()

	be.bufferByName[name] = buf.Bytes()

	return nil
}
