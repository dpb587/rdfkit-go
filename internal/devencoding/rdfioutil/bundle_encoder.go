package rdfioutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/rdfio"
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

func (be *BundleEncoder) PutBundle(name string, statements rdfio.StatementList) error {
	ctx := context.Background()
	buf := &bytes.Buffer{}

	e := NewEncoder(buf, EncoderOptions{})

	slices.SortStableFunc(statements, CompareStatementsDeterministic)

	for _, s := range statements {
		err := e.PutStatement(ctx, s)
		if err != nil {
			return err
		}
	}

	e.Close()

	be.bufferByName[name] = buf.Bytes()

	return nil
}
