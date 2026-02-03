package html

import (
	"fmt"
	"io"

	"github.com/dpb587/inspecthtml-go/inspecthtml"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"golang.org/x/net/html"
)

type DocumentOption interface {
	apply(s *DocumentConfig)
	newDocument(r io.Reader) (*Document, error)
}

type Document struct {
	info          DocumentInfo
	root          *html.Node
	parseMetadata *inspecthtml.ParseMetadata
	nodesByID     map[string][]*html.Node
}

func ParseDocument(r io.Reader, opts ...DocumentOption) (*Document, error) {
	compiledOpts := &DocumentConfig{}

	for _, opt := range opts {
		opt.apply(compiledOpts)
	}

	return compiledOpts.newDocument(r)
}

func NewDocument(root *html.Node, location string) (*Document, error) {
	var locationURL *iriutil.ParsedIRI

	if len(location) > 0 {
		var err error

		locationURL, err = iriutil.ParseIRI(location)
		if err != nil {
			return nil, fmt.Errorf("parse location: %v", err)
		}
	}

	d := &Document{
		info: DocumentInfo{
			Location: location,
		},
		root: root,
	}

	if baseHref, ok := findFirstBaseHref(d.root); ok {
		baseURL, err := iriutil.ParseIRI(baseHref)
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

func (d *Document) GetInfo() DocumentInfo {
	return d.info
}

func (d *Document) GetRoot() *html.Node {
	return d.root
}

func (d *Document) GetNodeMetadata(n *html.Node) (*inspecthtml.NodeMetadata, bool) {
	if d.parseMetadata == nil {
		return nil, false
	}

	return d.parseMetadata.GetNodeMetadata(n)
}

func (d *Document) GetNodesByID(id string) []*html.Node {
	if d.nodesByID == nil {
		d.nodesByID = map[string][]*html.Node{}
		d.indexNodesById(d.root)
	}

	return d.nodesByID[id]
}

func (d *Document) indexNodesById(n *html.Node) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Namespace == "" && attr.Key == "id" {
				d.nodesByID[attr.Val] = append(d.nodesByID[attr.Val], n)

				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		d.indexNodesById(c)
	}
}
