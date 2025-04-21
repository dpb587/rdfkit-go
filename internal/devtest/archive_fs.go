package devtest

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
)

func OpenArchiveTarGz(fp string, fpMap func(v string) string) (map[string][]byte, error) {
	fh, err := os.OpenFile(fp, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	defer fh.Close()

	fhGzip, err := gzip.NewReader(fh)
	if err != nil {
		return nil, fmt.Errorf("gunzip: %v", err)
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

			return nil, fmt.Errorf("next: %v", err)
		}

		data, err := io.ReadAll(fhTar)
		if err != nil {
			return nil, fmt.Errorf("read: %v", err)
		}

		entriesKey := header.Name

		if fpMap != nil {
			entriesKey = fpMap(header.Name)
		}

		entries[entriesKey] = data
	}

	return entries, nil
}
