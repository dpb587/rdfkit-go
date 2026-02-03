package html

import (
	"fmt"
	"io"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspecthtml-go/inspecthtml"
	"github.com/dpb587/rdfkit-go/iri"
	"golang.org/x/net/html"
)

type DocumentConfig struct {
	location *string

	captureTextOffsets *bool
	initialTextOffset  *cursorio.TextOffset
}

var _ DocumentOption = DocumentConfig{}

func (o DocumentConfig) SetLocation(v string) DocumentConfig {
	o.location = &v

	return o
}

func (o DocumentConfig) SetCaptureTextOffsets(v bool) DocumentConfig {
	o.captureTextOffsets = &v

	return o
}

func (o DocumentConfig) SetInitialTextOffset(v cursorio.TextOffset) DocumentConfig {
	t := true

	o.captureTextOffsets = &t
	o.initialTextOffset = &v

	return o
}

func (b DocumentConfig) apply(s *DocumentConfig) {
	if b.location != nil {
		s.location = b.location
	}

	if b.captureTextOffsets != nil {
		s.captureTextOffsets = b.captureTextOffsets
	}

	if b.initialTextOffset != nil {
		s.initialTextOffset = b.initialTextOffset
	}
}

func (b DocumentConfig) newDocument(r io.Reader) (*Document, error) {
	var location string
	var locationURL *iri.ParsedIRI

	if b.location != nil {
		var err error

		location = *b.location
		locationURL, err = iri.ParseIRI(location)
		if err != nil {
			return nil, fmt.Errorf("parse location: %v", err)
		}
	}

	d := &Document{
		info: DocumentInfo{
			Location: location,
			BaseURL:  location,
		},
	}

	var err error

	if b.captureTextOffsets != nil && *b.captureTextOffsets {
		hopt := []inspecthtml.ParserOption{}

		if b.initialTextOffset != nil {
			hopt = append(hopt, inspecthtml.ParserConfig{}.SetInitialOffset(*b.initialTextOffset))
		}

		d.root, d.parseMetadata, err = inspecthtml.NewParser(r, hopt...).Parse()

		d.info.HasNodeMetadata = true
	} else {
		d.root, err = html.Parse(r)
	}

	if err != nil {
		return nil, err
	}

	if baseHref, ok := findFirstBaseHref(d.root); ok {
		baseURL, err := iri.ParseIRI(baseHref)
		if err != nil {
			// TODO warn
		} else {
			if baseURL.IsAbs() {
				d.info.BaseURL = baseURL.String()
			} else if locationURL != nil {
				d.info.BaseURL = locationURL.ResolveReference(baseURL).String()
			} else {
				d.info.BaseURL = baseHref
			}
		}
	}

	return d, nil
}
