package testingarchive

import (
	"bytes"
	"maps"
	"slices"
	"testing"
)

type Archive struct {
	entries map[string][]byte
}

func (a Archive) ListFiles() []string {
	return slices.Collect(maps.Keys(a.entries))
}

func (a Archive) HasFile(name string) bool {
	_, ok := a.entries[name]

	return ok
}

func (a Archive) GetFileBytes(t *testing.T, name string) []byte {
	data, ok := a.entries[name]
	if !ok {
		t.Fatalf("archive file not found: %s", name)
	}

	return data
}

func (a Archive) NewFileByteReader(t *testing.T, name string) *bytes.Reader {
	data, ok := a.entries[name]
	if !ok {
		t.Fatalf("archive file not found: %s", name)
	}

	return bytes.NewReader(data)
}
