package rdfa

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func (v *Decoder) xmlRender(n *html.Node, xmlnsKnown map[string]string) (string, error) {
	switch n.Namespace {
	case "":
		//
	case "math":
		xmlnsKnown[""] = "http://www.w3.org/1998/Math/MathML"
	case "svg":
		xmlnsKnown[""] = "http://www.w3.org/2000/svg"
	}

	for _, attr := range n.Attr {
		if attr.Namespace != "" {
			continue
		} else if attr.Key == "xmlns" {
			xmlnsKnown[""] = attr.Val
		} else if strings.HasPrefix(attr.Key, "xmlns:") {
			xmlnsKnown[attr.Key[6:]] = attr.Val
		}
	}

	buf := &bytes.Buffer{}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		rebuilt, xmlnsHasRoot, xmlnsMissing := v.xmlRebuild(c, map[string]string{})

		if _, xmlnsKnownRoot := xmlnsKnown[""]; !xmlnsHasRoot && xmlnsKnownRoot {
			rebuilt.Attr = append([]html.Attribute{
				{
					Key: "xmlns",
					Val: xmlnsKnown[""],
				},
			}, rebuilt.Attr...)
		}

		for k := range xmlnsMissing {
			if schema, known := xmlnsKnown[k]; known {
				rebuilt.Attr = append(rebuilt.Attr, html.Attribute{
					Key: "xmlns:" + k,
					Val: schema,
				})
			}
		}

		err := html.Render(buf, rebuilt)
		if err != nil {
			return "", fmt.Errorf("render: %v", err)
		}
	}

	raw := buf.String()

	// TODO less hacky
	re := regexp.MustCompile(`<([^/][^\s]*)(|\s+[^>]+)></([^>]+)>`)
	raw = re.ReplaceAllStringFunc(raw, func(match string) string {
		m := re.FindStringSubmatch(match)

		if m[1] == m[3] {
			return "<" + m[1] + m[2] + "/>"
		}

		return match
	})

	return raw, nil
}

func (v *Decoder) xmlRebuild(n *html.Node, xmlnsKnown map[string]string) (*html.Node, bool, map[string]struct{}) {
	xmlnsFound := map[string]struct{}{}

	nextNode := &html.Node{
		Type:      n.Type,
		DataAtom:  n.DataAtom,
		Data:      n.Data,
		Namespace: n.Namespace,
	}

	if n.DataAtom == 0x0 {
		keySplit := strings.SplitN(n.Data, ":", 2)
		if len(keySplit) == 2 {
			xmlnsFound[keySplit[0]] = struct{}{}
		}
	}

	for _, attr := range n.Attr {
		if attr.Namespace == "" && attr.Key == "data-turple-offset" {
			continue
		}

		if attr.Namespace == "" && attr.Key == "xmlns" {
			xmlnsKnown[""] = attr.Val
			xmlnsFound[""] = struct{}{}
		} else if strings.HasPrefix(attr.Key, "xmlns:") {
			xmlnsKnown[attr.Key[6:]] = attr.Val
			xmlnsFound[attr.Key[6:]] = struct{}{}
		} else {
			keySplit := strings.SplitN(attr.Key, ":", 2)
			if len(keySplit) == 2 {
				xmlnsFound[keySplit[0]] = struct{}{}
			}
		}

		nextNode.Attr = append(nextNode.Attr, attr)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nextChild, _, xmlnsMissingChild := v.xmlRebuild(c, xmlnsKnown)
		nextNode.AppendChild(nextChild)

		for k := range xmlnsMissingChild {
			xmlnsFound[k] = struct{}{}
		}
	}

	xmlnsMissing := map[string]struct{}{}

	for k := range xmlnsFound {
		if _, known := xmlnsKnown[k]; !known {
			xmlnsMissing[k] = struct{}{}
		}
	}

	_, xmlnsHasRoot := xmlnsFound[""]

	return nextNode, xmlnsHasRoot, xmlnsMissing
}
