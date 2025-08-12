package htmlmicrodata

import (
	"fmt"

	"golang.org/x/net/html"
)

type DecoderError_LaxContentAttribute struct {
	Node      *html.Node
	Attribute html.Attribute
}

var _ error = DecoderError_LaxContentAttribute{}

func (e DecoderError_LaxContentAttribute) Error() string {
	return fmt.Sprintf("node %q references attribute %q", e.Node.DataAtom, "content")
}
