package rdfdescriptionstruct

import (
	"errors"
	"strings"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
)

// defaultPrefixMap is used to expand compact IRIs in struct tags.
var defaultPrefixMap = iriutil.NewPrefixMap(rdfacontext.WidelyUsedInitialContext()...)

// tagInfo contains parsed information from a struct field's rdf tag.
type tagInfo struct {
	// kind is the tag type: "s" for subject, "o" for object
	kind string

	// predicate is the predicate IRI for object tags
	predicate rdf.IRI
}

// parseTag parses an rdf struct tag.
// Valid formats:
// - "s" - subject
// - "o,p={PREDICATE}" - object scalar (PREDICATE can be full IRI or compact IRI)
func parseTag(tag string, prefixes iriutil.PrefixMap) (*tagInfo, error) {
	if tag == "" {
		return nil, nil
	}

	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return nil, nil
	}

	switch parts[0] {
	case "s":
		if len(parts) != 1 {
			return nil, errors.New("subject tag must be 's' only")
		}
		return &tagInfo{kind: "s"}, nil

	case "o":
		if len(parts) < 2 {
			return nil, errors.New("object tag must have predicate in format 'p={PREDICATE}'")
		}

		predicatePart := parts[1]
		if !strings.HasPrefix(predicatePart, "p=") {
			return nil, errors.New("object tag predicate must be in format 'p={PREDICATE}'")
		}

		predicateStr := predicatePart[2:] // Remove "p=" prefix
		if predicateStr == "" {
			return nil, errors.New("object tag predicate cannot be empty")
		}

		predicate := expandPredicate(predicateStr, prefixes)

		info := &tagInfo{
			kind:      "o",
			predicate: predicate,
		}

		// Check for unknown parameters
		if len(parts) > 2 {
			return nil, errors.New("unknown object tag parameter: " + parts[2])
		}

		return info, nil

	default:
		return nil, errors.New("tag must start with 's' or 'o'")
	}
}

// expandPredicate expands a compact IRI to a full IRI if necessary.
// If the predicate contains a colon, it's treated as a compact IRI (prefix:reference).
// Otherwise, it's returned as-is (assumed to be a full IRI).
func expandPredicate(predicate string, prefixes iriutil.PrefixMap) rdf.IRI {
	// Check if it's a compact IRI (contains colon but not at start)
	if colonIdx := strings.Index(predicate, ":"); colonIdx > 0 {
		prefix := predicate[:colonIdx]
		reference := predicate[colonIdx+1:]

		// Try to expand using the provided prefix map
		if expanded, ok := prefixes.ExpandPrefix(prefix, reference); ok {
			return expanded
		}
	}

	// Return as-is (either a full IRI or unrecognized compact form)
	return rdf.IRI(predicate)
}
