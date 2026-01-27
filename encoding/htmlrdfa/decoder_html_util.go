package htmlrdfa

import (
	"bytes"
	"fmt"

	"github.com/dpb587/rdfkit-go/rdf"
	"golang.org/x/net/html"
)

// https://www.w3.org/2011/rdfa-context/rdfa-1.1
var w3_2011_rdfacontext_rdfa11_TermMappings = map[string]rdf.IRI{
	"describedby": "http://www.w3.org/2007/05/powder-s#describedby",
	"license":     "http://www.w3.org/1999/xhtml/vocab#license",
	"role":        "http://www.w3.org/1999/xhtml/vocab#role",
}

// https://www.w3.org/2011/rdfa-context/xhtml-rdfa-1.1
var w3_2011_rdfacontext_xhtmlrdfa11_TermMappings = map[string]rdf.IRI{
	// XHTML Vocabulary (https://www.w3.org/1999/xhtml/vocab/)
	"alternate":  "http://www.w3.org/1999/xhtml/vocab#alternate",
	"appendix":   "http://www.w3.org/1999/xhtml/vocab#appendix",
	"cite":       "http://www.w3.org/1999/xhtml/vocab#cite",
	"bookmark":   "http://www.w3.org/1999/xhtml/vocab#bookmark",
	"contents":   "http://www.w3.org/1999/xhtml/vocab#contents",
	"chapter":    "http://www.w3.org/1999/xhtml/vocab#chapter",
	"copyright":  "http://www.w3.org/1999/xhtml/vocab#copyright",
	"first":      "http://www.w3.org/1999/xhtml/vocab#first",
	"glossary":   "http://www.w3.org/1999/xhtml/vocab#glossary",
	"help":       "http://www.w3.org/1999/xhtml/vocab#help",
	"icon":       "http://www.w3.org/1999/xhtml/vocab#icon",
	"index":      "http://www.w3.org/1999/xhtml/vocab#index",
	"last":       "http://www.w3.org/1999/xhtml/vocab#last",
	"license":    "http://www.w3.org/1999/xhtml/vocab#license",
	"meta":       "http://www.w3.org/1999/xhtml/vocab#meta",
	"next":       "http://www.w3.org/1999/xhtml/vocab#next",
	"prev":       "http://www.w3.org/1999/xhtml/vocab#prev",
	"previous":   "http://www.w3.org/1999/xhtml/vocab#previous",
	"section":    "http://www.w3.org/1999/xhtml/vocab#section",
	"start":      "http://www.w3.org/1999/xhtml/vocab#start",
	"stylesheet": "http://www.w3.org/1999/xhtml/vocab#stylesheet",
	"subsection": "http://www.w3.org/1999/xhtml/vocab#subsection",
	"top":        "http://www.w3.org/1999/xhtml/vocab#top",
	"up":         "http://www.w3.org/1999/xhtml/vocab#up",

	// W3C Recommendations (https://www.w3.org/TR/2002/REC-P3P-20020416/)
	"p3pv1": "http://www.w3.org/1999/xhtml/vocab#p3pv1",
}

// https://html.spec.whatwg.org/multipage/links.html#linkTypes
// haven't found official commentary about this, but ignoring node+link types which:
// (1) are listed by whatwg
// (2) not listed as a Hyperlink in the table
// (3) noted for "keyword does not create a hyperlink" in details
// reviewed 2024-08-19
var htmlIgnoredLinkRels = map[string]struct{}{
	"alternate/form":        {}, // not allowed
	"canonical/a":           {}, // not allowed
	"canonical/area":        {}, // not allowed
	"canonical/form":        {}, // not allowed
	"author/form":           {}, // not allowed
	"bookmark/link":         {}, // not allowed
	"bookmark/form":         {}, // not allowed
	"dns-prefetch/link":     {}, // External Resource
	"dns-prefetch/a":        {}, // not allowed
	"dns-prefetch/area":     {}, // not allowed
	"dns-prefetch/form":     {}, // not allowed
	"expect/link":           {}, // Internal Resource
	"expect/a":              {}, // not allowed
	"expect/area":           {}, // not allowed
	"expect/form":           {}, // not allowed
	"external/link":         {}, // not allowed
	"external/a":            {}, // Annotation
	"external/area":         {}, // Annotation
	"external/form":         {}, // Annotation
	"icon/link":             {}, // External Resource
	"icon/a":                {}, // not allowed
	"icon/area":             {}, // not allowed
	"icon/form":             {}, // not allowed
	"manifest/link":         {}, // External Resource
	"manifest/a":            {}, // not allowed
	"manifest/area":         {}, // not allowed
	"manifest/form":         {}, // not allowed
	"modulepreload/link":    {}, // External Resource
	"modulepreload/a":       {}, // not allowed
	"modulepreload/area":    {}, // not allowed
	"modulepreload/form":    {}, // not allowed
	"nofollow/link":         {}, // not allowed
	"nofollow/a":            {}, // Annotation
	"nofollow/area":         {}, // Annotation
	"nofollow/form":         {}, // Annotation
	"noopener/link":         {}, // not allowed
	"noopener/a":            {}, // Annotation
	"noopener/area":         {}, // Annotation
	"noopener/form":         {}, // Annotation
	"noreferrer/link":       {}, // not allowed
	"noreferrer/a":          {}, // Annotation
	"noreferrer/area":       {}, // Annotation
	"noreferrer/form":       {}, // Annotation
	"opener/link":           {}, // not allowed
	"opener/a":              {}, // Annotation
	"opener/area":           {}, // Annotation
	"opener/form":           {}, // Annotation
	"pingback/link":         {}, // External Resource
	"pingback/a":            {}, // not allowed
	"pingback/area":         {}, // not allowed
	"pingback/form":         {}, // not allowed
	"preconnect/link":       {}, // External Resource
	"preconnect/a":          {}, // not allowed
	"preconnect/area":       {}, // not allowed
	"preconnect/form":       {}, // not allowed
	"prefetch/link":         {}, // External Resource
	"prefetch/a":            {}, // not allowed
	"prefetch/area":         {}, // not allowed
	"prefetch/form":         {}, // not allowed
	"preload/link":          {}, // External Resource
	"preload/a":             {}, // not allowed
	"preload/area":          {}, // not allowed
	"preload/form":          {}, // not allowed
	"privacy-policy/form":   {}, // not allowed
	"stylesheet/link":       {}, // External Resource
	"stylesheet/a":          {}, // not allowed
	"stylesheet/area":       {}, // not allowed
	"stylesheet/form":       {}, // not allowed
	"tag/link":              {}, // not allowed
	"tag/form":              {}, // not allowed
	"terms-of-service/form": {}, // not allowed
}

func (v *Decoder) htmlRender(n *html.Node) (string, error) {
	buf := &bytes.Buffer{}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		rebuilt := v.htmlRebuild(c)

		err := html.Render(buf, rebuilt)
		if err != nil {
			return "", fmt.Errorf("render: %v", err)
		}
	}

	raw := buf.String()

	return raw, nil
}

func (v *Decoder) htmlRebuild(n *html.Node) *html.Node {
	nextNode := &html.Node{
		Type:      n.Type,
		DataAtom:  n.DataAtom,
		Data:      n.Data,
		Namespace: n.Namespace,
	}

	for _, attr := range n.Attr {
		if attr.Namespace == "" && attr.Key == "data-turple-offset" {
			continue
		}

		nextNode.Attr = append(nextNode.Attr, attr)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nextChild := v.htmlRebuild(c)
		nextNode.AppendChild(nextChild)
	}

	return nextNode
}
