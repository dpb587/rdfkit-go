package html

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// https://html.spec.whatwg.org/multipage/urls-and-fetching.html#document-base-url
func findFirstBaseHref(n *html.Node) (string, bool) {
	if n.DataAtom == atom.Base {
		for _, attr := range n.Attr {
			if attr.Namespace == "" && attr.Key == "href" {
				return attr.Val, true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if base, ok := findFirstBaseHref(c); ok {
			return base, true
		}
	}

	return "", false
}
