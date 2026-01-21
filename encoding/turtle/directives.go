package turtle

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

// DirectiveMode configures how base and prefix directives are written
type DirectiveMode int

const (
	// DirectiveMode_At uses `@base` and `@prefix` (default).
	DirectiveMode_At DirectiveMode = iota

	// DirectiveMode_SPARQL uses `BASE` and `PREFIX`.
	DirectiveMode_SPARQL

	// DirectiveMode_Disabled disables directive output.
	DirectiveMode_Disabled
)

type WriteDirectivesOptions struct {
	Base     string
	BaseMode DirectiveMode

	Prefixes   iriutil.PrefixMappingList
	PrefixMode DirectiveMode
}

func WriteDirectives(w io.Writer, opts WriteDirectivesOptions) (int64, error) {
	buf := &bytes.Buffer{}

	if len(opts.Base) > 0 && opts.BaseMode != DirectiveMode_Disabled {
		switch opts.BaseMode {
		case DirectiveMode_SPARQL:
			fmt.Fprintf(buf, "BASE <%s>\n", opts.Base)
		case DirectiveMode_At:
			fmt.Fprintf(buf, "@base <%s> .\n", opts.Base)
		default:
			return 0, fmt.Errorf("unknown base mode: %v", opts.BaseMode)
		}
	}

	if opts.PrefixMode != DirectiveMode_Disabled {
		for _, t := range opts.Prefixes {
			switch opts.PrefixMode {
			case DirectiveMode_SPARQL:
				fmt.Fprintf(buf, "PREFIX %s: <%s>\n", t.Prefix, t.Expanded)
			case DirectiveMode_At:
				fmt.Fprintf(buf, "@prefix %s: <%s> .\n", t.Prefix, t.Expanded)
			default:
				return 0, fmt.Errorf("unknown prefix mode: %v", opts.PrefixMode)
			}
		}
	}

	return buf.WriteTo(w)
}
