package iri

import (
	"net/url"
	"strings"
	_ "unsafe"
)

// ParsedIRI is a light wrapper to url.URL with the following notable exceptions:
//
//   - raw path and raw fragments are preserved
//   - empty fragment data is always included, if included in the original string
type ParsedIRI struct {
	u             *url.URL
	forceFragment bool
	isOpaque      bool // true if this is a non-hierarchical URI (uses Opaque field)
}

func ParseIRI(s string) (*ParsedIRI, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	isOpaque := false

	// For non-hierarchical schemes (not http/https), use Opaque to avoid '//' authority
	if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "file" && u.Host == "" && u.Opaque == "" {
		isOpaque = true
		// url.Parse puts everything after scheme: into Path with leading '/'
		// Move it to Opaque without the leading '/'
		if u.Path != "" {
			if u.Path[0] == '/' {
				u.Opaque = u.Path[1:]
			} else {
				u.Opaque = u.Path
			}
			u.Path = ""
			u.RawPath = ""
		} else {
			// Empty path - set Opaque to empty string to mark as non-hierarchical
			u.Opaque = ""
		}
	}

	return &ParsedIRI{
		u:             u,
		forceFragment: strings.HasSuffix(s, "#"),
		isOpaque:      isOpaque,
	}, nil
}

func (u *ParsedIRI) URL() *url.URL {
	uu := *u.u

	return &uu
}

func (u *ParsedIRI) IsAbs() bool {
	return u.u.IsAbs()
}

func (u *ParsedIRI) DropFragment() {
	u.forceFragment = false
	u.u.Fragment = ""
	u.u.RawFragment = ""
}

func (u *ParsedIRI) Parse(ref string) (*ParsedIRI, error) {
	// duplicated from stdlib

	refIRI, err := ParseIRI(ref)
	if err != nil {
		return nil, err
	}

	return u.ResolveReference(refIRI), nil
}

// bad; marginally better than duplicating private encode(*, encodePath) behavior?
//
//go:linkname badSetPath net/url.(*URL).setPath
func badSetPath(u *url.URL, path string)

func (iri *ParsedIRI) ResolveReference(ref *ParsedIRI) *ParsedIRI {
	u, url := iri.u, *ref.u
	uPath, refuPath := u.EscapedPath(), ref.u.EscapedPath()

	if len(iri.u.RawPath) > 0 {
		uPath = iri.u.RawPath
	}

	if len(ref.u.RawPath) > 0 {
		refuPath = ref.u.RawPath
	}

	forceFragment := iri.forceFragment || ref.forceFragment

	// [dpb] below mostly duplicated from stdlib

	if ref.u.Scheme == "" {
		url.Scheme = u.Scheme
	}
	if ref.u.Scheme != "" || ref.u.Host != "" || ref.u.User != nil {
		// The "absoluteURI" or "net_path" cases.
		// We can ignore the error from setPath since we know we provided a
		// validly-escaped path.
		badSetPath(&url, resolvePath(refuPath, ""))
		return &ParsedIRI{
			u:             &url,
			forceFragment: forceFragment,
			isOpaque:      ref.isOpaque,
		}
	}
	if ref.u.Opaque != "" {
		url.User = nil
		url.Host = ""
		url.Path = ""
		return &ParsedIRI{
			u:             &url,
			forceFragment: forceFragment,
			isOpaque:      true,
		}
	}
	if ref.u.Path == "" && !ref.u.ForceQuery && ref.u.RawQuery == "" {
		url.RawQuery = u.RawQuery
		if ref.u.Fragment == "" {
			url.Fragment = u.Fragment
			url.RawFragment = u.RawFragment
		}
	}
	if ref.u.Path == "" && u.Opaque != "" {
		url.Opaque = u.Opaque
		url.User = nil
		url.Host = ""
		url.Path = ""
		return &ParsedIRI{
			u:             &url,
			forceFragment: forceFragment,
			isOpaque:      true,
		}
	}

	// Handle resolving relative refs against opaque (non-hierarchical) base URIs
	if u.Opaque != "" || iri.isOpaque {
		// For opaque URIs, resolve the reference path-like
		var resolved string
		if refuPath != "" {
			if refuPath[0] == '/' {
				// Absolute path reference
				resolved = refuPath
			} else {
				// Relative path reference - merge with opaque part
				base := u.Opaque
				i := strings.LastIndex(base, "/")
				if i >= 0 {
					resolved = resolvePath(base[:i+1]+refuPath, "")
				} else {
					resolved = resolvePath(refuPath, "")
				}
				// Remove leading '/' from resolvePath result for opaque
				if len(resolved) > 0 && resolved[0] == '/' {
					resolved = resolved[1:]
				}
			}
		} else {
			resolved = u.Opaque
		}
		url.Opaque = resolved
		url.User = nil
		url.Host = ""
		url.Path = ""
		url.RawPath = ""
		return &ParsedIRI{
			u:             &url,
			forceFragment: forceFragment,
			isOpaque:      true,
		}
	}

	// The "abs_path" or "rel_path" cases.
	url.Host = u.Host
	url.User = u.User

	if uPath == "" {
		// [dpb] handle empty base with relative ref - don't force absolute path

		if len(refuPath) > 0 {
			// force dot-segment resolution
			resolved := resolvePath(refuPath, "")

			if refuPath[0] != '/' {
				resolved = resolved[1:]
			}

			badSetPath(&url, resolved)
		}
	} else {
		badSetPath(&url, resolvePath(uPath, refuPath))
	}

	return &ParsedIRI{
		u:             &url,
		forceFragment: forceFragment,
		isOpaque:      false,
	}
}

func (iri *ParsedIRI) String() string {
	// hacky to strings-replace values?
	// std String() relies on private escape functions that would need to be duplicated

	s := iri.u.String()

	if len(iri.u.RawPath) > 0 {
		s = strings.Replace(s, iri.u.EscapedPath(), iri.u.RawPath, 1)
	}

	if len(iri.u.RawFragment) > 0 {
		s = strings.Replace(s, "#"+iri.u.EscapedFragment(), "#"+iri.u.RawFragment, 1)
	} else if iri.forceFragment && !strings.Contains(s, "#") {
		s += "#"
	}

	return s
}

// fully duplicated from stdlib
func resolvePath(base, ref string) string {
	var full string
	if ref == "" {
		full = base
	} else if ref[0] != '/' {
		i := strings.LastIndex(base, "/")
		full = base[:i+1] + ref
	} else {
		full = ref
	}
	if full == "" {
		return ""
	}

	var (
		elem string
		dst  strings.Builder
	)
	first := true
	remaining := full
	// We want to return a leading '/', so write it now.
	dst.WriteByte('/')
	found := true
	for found {
		elem, remaining, found = strings.Cut(remaining, "/")
		if elem == "." {
			first = false
			// drop
			continue
		}

		if elem == ".." {
			// Ignore the leading '/' we already wrote.
			str := dst.String()[1:]
			index := strings.LastIndexByte(str, '/')

			dst.Reset()
			dst.WriteByte('/')
			if index == -1 {
				first = true
			} else {
				dst.WriteString(str[:index])
			}
		} else {
			if !first {
				dst.WriteByte('/')
			}
			dst.WriteString(elem)
			first = false
		}
	}

	if elem == "." || elem == ".." {
		dst.WriteByte('/')
	}

	// We wrote an initial '/', but we don't want two.
	r := dst.String()
	if len(r) > 1 && r[1] == '/' {
		r = r[1:]
	}
	return r
}
