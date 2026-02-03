package curie

import "strings"

func ParseCURIE(v string) (CURIE, bool) {
	if len(v) == 0 {
		return CURIE{}, false
	}

	var res CURIE

	if v[0] == '[' && v[len(v)-1] == ']' {
		res.Safe = true
		v = v[1 : len(v)-1]
	}

	// TODO NCName
	parts := strings.SplitN(v, ":", 2)
	if len(parts) == 2 {
		res.Prefix = parts[0]
		res.Reference = parts[1]
	} else {
		res.DefaultPrefix = true
		res.Reference = parts[0]
	}

	return res, true
}

func ParseCURIEs(v string) CURIEs {
	var res CURIEs

	for _, vv := range strings.Fields(v) {
		if len(vv) == 0 {
			continue
		}

		c, ok := ParseCURIE(vv)
		if !ok {
			continue
		}

		res = append(res, c)
	}

	return res
}
