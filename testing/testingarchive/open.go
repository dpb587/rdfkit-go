package testingarchive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"testing"
)

func OpenTarGz(t *testing.T, fp string, fpMap func(v string) string) Archive {
	fh, err := os.OpenFile(fp, os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	defer fh.Close()

	fhGzip, err := gzip.NewReader(fh)
	if err != nil {
		t.Fatalf("gunzip: %v", err)
	}

	defer fhGzip.Close()

	fhTar := tar.NewReader(fhGzip)

	entries := make(map[string][]byte)

	for {
		header, err := fhTar.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			t.Fatalf("next: %v", err)
		}

		data, err := io.ReadAll(fhTar)
		if err != nil {
			t.Fatalf("read: %v", err)
		}

		entriesKey := header.Name

		if fpMap != nil {
			entriesKey = fpMap(header.Name)
		}

		entries[entriesKey] = data
	}

	return Archive{entries: entries}
}
