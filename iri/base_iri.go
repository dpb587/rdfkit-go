package iri

import (
	"strings"
)

type BaseIRI struct {
	original string
	parsed   *ParsedIRI

	rootIndex      int
	directoryIndex int
	resourceIndex  int
	queryIndex     int
	fragmentIndex  int
}

func ParseBaseIRI(v string) (*BaseIRI, error) {
	parsed, err := ParseIRI(v)
	if err != nil {
		return nil, err
	}

	return NewBaseIRI(parsed), nil
}

func NewBaseIRI(parsed *ParsedIRI) *BaseIRI {
	v := parsed.String()

	baseFragment, _, _ := strings.Cut(v, "#")
	baseQuery, _, _ := strings.Cut(baseFragment, "?")

	rb := &BaseIRI{
		original:      v,
		parsed:        parsed,
		resourceIndex: len(baseQuery),
	}

	if baseFragment != v {
		rb.fragmentIndex = len(baseFragment)
	} else {
		rb.fragmentIndex = -1
	}

	if baseQuery != baseFragment {
		rb.queryIndex = len(baseQuery)
	} else {
		rb.queryIndex = -1
	}

	if parsed.IsAbs() {
		baseRoot, _ := parsed.Parse("/")
		baseSubpathParsed, _ := parsed.Parse("./")
		rb.rootIndex = len(baseRoot.String())
		rb.directoryIndex = len(baseSubpathParsed.String())
	} else {
		rb.rootIndex = -1
		rb.directoryIndex = -1
	}

	return rb
}

func (rb *BaseIRI) IsAbs() bool {
	return rb.rootIndex != -1
}

func (rb *BaseIRI) String() string {
	return rb.original
}

func (rb *BaseIRI) Parse(v string) (*ParsedIRI, error) {
	return rb.parsed.Parse(v)
}

func (rb *BaseIRI) ResolveReference(ref *ParsedIRI) *ParsedIRI {
	return rb.parsed.ResolveReference(ref)
}

func (rb *BaseIRI) RelativizeIRI(v string) (string, bool) {
	if len(v) > len(rb.original) {
		if rb.fragmentIndex == -1 && v[len(rb.original)] == '#' {
			return v[len(rb.original):], true
		} else if rb.queryIndex == -1 && v[len(rb.original)] == '?' {
			return v[len(rb.original):], true
		}
	}

	if rb.rootIndex == -1 {
		return "", false
	} else if len(v) < rb.rootIndex || rb.original[0:rb.rootIndex] != v[:rb.rootIndex] {
		return "", false
	} else if rb.original == v {
		return "", true
	}

	if len(v) > rb.resourceIndex {
		switch v[rb.resourceIndex] {
		case '#':
			// dropping query
			return v[rb.directoryIndex:], true
		case '?':
			return v[rb.resourceIndex:], true
		}
	}

	if len(v) >= rb.directoryIndex && rb.original[0:rb.directoryIndex] == v[:rb.directoryIndex] {
		return v[rb.directoryIndex:], true
	}

	return v[rb.rootIndex-1:], true
}
