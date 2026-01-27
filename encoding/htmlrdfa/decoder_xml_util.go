package htmlrdfa

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func (v *Decoder) xmlRender(n *html.Node) (string, error) {
	xmlnsKnown := make(map[string]string)

	for p := n; p != nil; p = p.Parent {
		if _, exists := xmlnsKnown[""]; !exists {
			switch p.Namespace {
			case "math":
				xmlnsKnown[""] = "http://www.w3.org/1998/Math/MathML"
			case "svg":
				xmlnsKnown[""] = "http://www.w3.org/2000/svg"
			}
		}

		for _, attr := range p.Attr {
			if attr.Namespace != "" {
				continue
			} else if attr.Key == "xmlns" {
				if _, exists := xmlnsKnown[""]; !exists {
					xmlnsKnown[""] = attr.Val
				}
			} else if strings.HasPrefix(attr.Key, "xmlns:") {
				prefix := attr.Key[6:]
				if _, exists := xmlnsKnown[prefix]; !exists {
					xmlnsKnown[prefix] = attr.Val
				}
			} else if attr.Key == "prefix" {
				// rdfa propagated as xmlns
				fields := strings.Fields(strings.TrimSpace(attr.Val))

				for fieldIdx := 0; fieldIdx+1 < len(fields); fieldIdx += 2 {
					prefixTerm := strings.ToLower(fields[fieldIdx])
					if strings.HasSuffix(prefixTerm, ":") {
						prefix := prefixTerm[:len(prefixTerm)-1]
						if _, exists := xmlnsKnown[prefix]; !exists {
							xmlnsKnown[prefix] = fields[fieldIdx+1]
						}
					}
				}
			}
		}
	}

	buf := &bytes.Buffer{}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		rebuilt, xmlnsHasRoot, _ := v.xmlRebuild(c, map[string]string{})
		var attrModified bool

		for k, schema := range xmlnsKnown {
			if len(k) == 0 && !xmlnsHasRoot {
				rebuilt.Attr = append(rebuilt.Attr, html.Attribute{
					Key: "xmlns",
					Val: schema,
				})

				continue
			}

			rebuilt.Attr = append(rebuilt.Attr, html.Attribute{
				Key: "xmlns:" + k,
				Val: schema,
			})

			attrModified = true
		}

		if attrModified {
			// not currently trying to do a full, recursive canonicalization
			slices.SortStableFunc(rebuilt.Attr, v.xmlExtc14n)
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

func (*Decoder) xmlExtc14n(a, b html.Attribute) int {
	aIsDefaultNS := a.Namespace == "" && a.Key == "xmlns"
	bIsDefaultNS := b.Namespace == "" && b.Key == "xmlns"
	if aIsDefaultNS != bIsDefaultNS {
		if aIsDefaultNS {
			return -1
		}

		return 1
	}

	aIsNSDecl := a.Namespace == "" && strings.HasPrefix(a.Key, "xmlns:")
	bIsNSDecl := b.Namespace == "" && strings.HasPrefix(b.Key, "xmlns:")
	if aIsNSDecl && bIsNSDecl {
		return strings.Compare(a.Key, b.Key)
	} else if aIsNSDecl != bIsNSDecl {
		if aIsNSDecl {
			return -1
		}
		return 1
	}

	aIsUnqualified := a.Namespace == ""
	bIsUnqualified := b.Namespace == ""
	if aIsUnqualified && bIsUnqualified {
		return strings.Compare(a.Key, b.Key)
	} else if aIsUnqualified != bIsUnqualified {
		if aIsUnqualified {
			return -1
		}
		return 1
	}

	if a.Namespace != b.Namespace {
		return strings.Compare(a.Namespace, b.Namespace)
	}

	return strings.Compare(a.Key, b.Key)
}
