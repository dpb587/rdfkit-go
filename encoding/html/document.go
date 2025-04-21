package html

import (
	"io"

	"github.com/dpb587/inspecthtml-go/inspecthtml"
	"golang.org/x/net/html"
)

// DocumentOption is a functional option for [ParseDocument] to affect its behavior.
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
