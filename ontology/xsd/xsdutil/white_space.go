package xsdutil

import (
	"regexp"
	"strings"
)

// https://www.w3.org/TR/xmlschema-2/#rf-whiteSpace

// No normalization is done, the value is not changed (this is the behavior required by [XML 1.0 (Second Edition)] for element content)
func WhiteSpacePreserve(v string) string {
	return v
}

var whiteSpaceReplace = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ")

// All occurrences of #x9 (tab), #xA (line feed) and #xD (carriage return) are replaced with #x20 (space)
func WhiteSpaceReplace(v string) string {
	return whiteSpaceReplace.Replace(v)
}

var whiteSpaceCollapseRE = regexp.MustCompile(` +`)

// After the processing implied by replace, contiguous sequences of #x20's are collapsed to a single #x20, and leading and trailing #x20's are removed.
func WhiteSpaceCollapse(v string) string {
	return strings.TrimRight(strings.TrimLeft(whiteSpaceCollapseRE.ReplaceAllString(whiteSpaceReplace.Replace(v), " "), " "), " ")
}
