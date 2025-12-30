package rdfioutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
)

type BundleEncoder struct {
	w            io.Writer
	bufferByName map[string][]byte
}

func NewBundleEncoder(w io.Writer) *BundleEncoder {
	return &BundleEncoder{
		w:            w,
		bufferByName: make(map[string][]byte),
	}
}

func (be *BundleEncoder) Close() error {
	keys := slices.Collect(maps.Keys(be.bufferByName))
	slices.SortFunc(keys, strings.Compare)

	for nameIdx, name := range keys {
		if nameIdx > 0 {
			be.w.Write([]byte("\n"))
		}

		fmt.Fprintf(be.w, "=== %s\n", name)
		fmt.Fprintf(be.w, "\n")

		be.w.Write(be.bufferByName[name])
	}

	return nil
}

func (be *BundleEncoder) PutQuadsBundle(name string, statements encodingtest.QuadStatementList) error {
	ctx := context.Background()

	buf := &bytes.Buffer{}

	e := encodingtest.NewQuadsEncoder(buf, encodingtest.QuadsEncoderOptions{})

	for _, statement := range statements {
		err := e.AddQuadStatement(ctx, statement)
		if err != nil {
			return err
		}
	}

	e.Close()

	be.bufferByName[name] = buf.Bytes()

	return nil
}

func (be *BundleEncoder) PutTriplesBundle(name string, statements encodingtest.TripleStatementList) error {
	ctx := context.Background()

	buf := &bytes.Buffer{}

	e := encodingtest.NewTriplesEncoder(buf, encodingtest.TriplesEncoderOptions{})

	for _, statement := range statements {
		err := e.AddTripleStatement(ctx, statement)
		if err != nil {
			return err
		}
	}

	e.Close()

	be.bufferByName[name] = buf.Bytes()

	return nil
}
