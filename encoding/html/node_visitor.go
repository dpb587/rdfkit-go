package html

import "golang.org/x/net/html"

type NodeVisitor interface {
	VisitNode(*html.Node) error
}

//

type DocumentRootTransformerFunc func(*html.Node) error

var _ NodeVisitor = DocumentRootTransformerFunc(nil)

func (f DocumentRootTransformerFunc) VisitNode(n *html.Node) error {
	return f(n)
}
