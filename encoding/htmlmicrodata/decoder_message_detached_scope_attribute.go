package htmlmicrodata

import (
	"github.com/dpb587/inspecthtml-go/inspecthtml"
	"github.com/dpb587/rdfkit-go/encoding"
	"golang.org/x/net/html"
)

type DecoderMessage_DetachedScopeAttribute struct {
	Decoder *Decoder

	Node     *html.Node
	NodeAttr int

	AttrName string
}

var _ encoding.DecoderMessage = DecoderMessage_DetachedScopeAttribute{}

func (m DecoderMessage_DetachedScopeAttribute) GetDecoder() encoding.Decoder {
	return m.Decoder
}

func (m DecoderMessage_DetachedScopeAttribute) GetNodeAttrMetadata() *inspecthtml.NodeAttributeMetadata {
	metadata, ok := m.Decoder.doc.GetNodeMetadata(m.Node)
	if !ok {
		return nil
	}

	return metadata.TagAttr[m.NodeAttr]
}
